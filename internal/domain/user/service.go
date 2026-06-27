package user

import (
	"errors"
	"spotsync/internal/auth"
	"spotsync/internal/domain/user/dto"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetProfile(userID uint) (*dto.UserResponse, error)
	RefreshToken(userID uint, role string) (string, error)
}

type userService struct {
	repo       UserRepository
	jwtService auth.JWTService
}

func NewUserService(repo UserRepository, jwtService auth.JWTService) UserService {
	return &userService{
		repo:       repo,
		jwtService: jwtService,
	}
}

func (s *userService) Register(req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	// Check if email already exists
	existing, err := s.repo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password (bcrypt cost 10-12 according to specification)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		return nil, err
	}

	role := req.Role
	if role == "" {
		role = "driver"
	}

	user := &User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Generate tokens on registration
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Token:        accessToken,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *userService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens on login
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        accessToken,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *userService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (s *userService) RefreshToken(userID uint, role string) (string, error) {
	// Verify user still exists in the database
	_, err := s.repo.FindByID(userID)
	if err != nil {
		return "", err
	}

	// Generate new access token
	return s.jwtService.GenerateAccessToken(userID, role)
}
