package reservation

import (
	"errors"
	"spotsync/internal/domain/reservation/dto"
)

var (
	ErrForbidden            = errors.New("access forbidden: you can only cancel your own reservations")
	ErrReservationNotActive = errors.New("reservation is not active")
)

type ReservationService interface {
	ReserveSpot(userID uint, req *dto.CreateReservationRequest) (*dto.CreateReservationResponse, error)
	GetMyReservations(userID uint) ([]dto.MyReservationResponse, error)
	GetAllReservations() ([]dto.AdminReservationResponse, error)
	CancelReservation(userID uint, userRole string, reservationID uint) error
}

type reservationService struct {
	repo ReservationRepository
}

func NewReservationService(repo ReservationRepository) ReservationService {
	return &reservationService{repo: repo}
}

func (s *reservationService) ReserveSpot(userID uint, req *dto.CreateReservationRequest) (*dto.CreateReservationResponse, error) {
	res, err := s.repo.CreateWithLock(userID, req.ZoneID, req.LicensePlate)
	if err != nil {
		return nil, err
	}

	return &dto.CreateReservationResponse{
		ID:           res.ID,
		UserID:       res.UserID,
		ZoneID:       res.ZoneID,
		LicensePlate: res.LicensePlate,
		Status:       res.Status,
		CreatedAt:    res.CreatedAt,
		UpdatedAt:    res.UpdatedAt,
	}, nil
}

func (s *reservationService) GetMyReservations(userID uint) ([]dto.MyReservationResponse, error) {
	list, err := s.repo.FindMyReservations(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.MyReservationResponse
	for _, item := range list {
		res = append(res, dto.MyReservationResponse{
			ID:           item.ID,
			LicensePlate: item.LicensePlate,
			Status:       item.Status,
			Zone: dto.ReservationZoneResponse{
				ID:   item.Zone.ID,
				Name: item.Zone.Name,
				Type: item.Zone.Type,
			},
			CreatedAt:    item.CreatedAt,
		})
	}
	return res, nil
}

func (s *reservationService) GetAllReservations() ([]dto.AdminReservationResponse, error) {
	list, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.AdminReservationResponse
	for _, item := range list {
		res = append(res, dto.AdminReservationResponse{
			ID:     item.ID,
			UserID: item.UserID,
			User: dto.ReservationUserResponse{
				ID:    item.User.ID,
				Name:  item.User.Name,
				Email: item.User.Email,
			},
			ZoneID: item.ZoneID,
			Zone: dto.ReservationZoneResponse{
				ID:   item.Zone.ID,
				Name: item.Zone.Name,
				Type: item.Zone.Type,
			},
			LicensePlate: item.LicensePlate,
			Status:       item.Status,
			CreatedAt:    item.CreatedAt,
			UpdatedAt:    item.UpdatedAt,
		})
	}
	return res, nil
}

func (s *reservationService) CancelReservation(userID uint, userRole string, reservationID uint) error {
	res, err := s.repo.FindByID(reservationID)
	if err != nil {
		return err
	}

	// Drivers can only cancel their own reservations; admins can cancel any.
	if userRole != "admin" && res.UserID != userID {
		return ErrForbidden
	}

	if res.Status != "active" {
		return ErrReservationNotActive
	}

	return s.repo.UpdateStatus(reservationID, "cancelled")
}
