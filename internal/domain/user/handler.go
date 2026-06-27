package user

import (
	"net/http"
	"time"

	"spotsync/internal/domain/user/dto"

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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request payload",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Validation errors",
			"errors":  err.Error(),
		})
	}

	res, err := h.service.Register(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// Set access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(accessCookie)

	// Set refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(refreshCookie)

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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request payload",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Validation errors",
			"errors":  err.Error(),
		})
	}

	res, err := h.service.Login(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
		})
	}

	// Set access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    res.AccessToken,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(accessCookie)

	// Set refresh token cookie
	refreshCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(refreshCookie)

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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized: user ID not found in context",
		})
	}

	res, err := h.service.GetProfile(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "User profile not found",
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
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Unauthorized: invalid claims in context",
		})
	}

	token, err := h.service.RefreshToken(userID, role)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"success": false,
			"message": "Could not refresh token",
			"errors":  err.Error(),
		})
	}

	// Set updated access token cookie
	accessCookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(accessCookie)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Token refreshed successfully",
		"data": map[string]string{
			"token": token,
		},
	})
}
