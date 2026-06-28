package user

import (
	"net/http"
	"time"

	"spotsync/internal/domain/user/dto"
	"spotsync/internal/httpResponse"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register handles user registration.
func (h *UserHandler) Register(c echo.Context) error {
	var req dto.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Details: err.Error(),
		})
	}

	res, err := h.service.Register(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "User registered successfully",
		"data":    res,
	})
}

// Login handles user authentication.
func (h *UserHandler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Validation errors",
			Details: err.Error(),
		})
	}

	res, err := h.service.Login(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Set "access_token" cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    res.Token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(accessCookie)

	// Also set "token" cookie
	tokenCookie := &http.Cookie{
		Name:     "token",
		Value:    res.Token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(tokenCookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Login successful",
		"data":    res,
	})
}

// GetMe retrieves the current authenticated user's profile.
func (h *UserHandler) GetMe(c echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: user ID not found in context",
		})
	}

	res, err := h.service.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "User profile not found",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "User profile retrieved successfully",
		"data":    res,
	})
}

// RefreshToken generates a new JWT token using the current token's claims.
func (h *UserHandler) RefreshToken(c echo.Context) error {
	userIDVal := c.Get("user_id")
	roleVal := c.Get("role")

	userID, ok1 := userIDVal.(uint)
	role, ok2 := roleVal.(string)
	if !ok1 || !ok2 {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: invalid claims in context",
		})
	}

	token, err := h.service.RefreshToken(userID, role)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Could not refresh token",
			Details: err.Error(),
		})
	}

	// Set updated "access_token" cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(accessCookie)

	// Also set updated "token" cookie
	tokenCookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(tokenCookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Token refreshed successfully",
		"data": map[string]string{
			"token": token,
		},
	})
}
