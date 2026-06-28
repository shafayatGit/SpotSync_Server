package reservation

import (
	"errors"
	"net/http"
	"strconv"

	"spotsync/internal/domain/reservation/dto"
	"spotsync/internal/httpResponse"

	"github.com/labstack/echo/v4"
)

type ReservationHandler struct {
	service ReservationService
}

func NewReservationHandler(service ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) ReserveSpot(c echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: user ID not found in context",
		})
	}

	var req dto.CreateReservationRequest
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

	res, err := h.service.ReserveSpot(userID, &req)
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ErrZoneNotFound) {
			code = http.StatusNotFound
		} else if errors.Is(err, ErrZoneFull) {
			code = http.StatusConflict
		}

		return c.JSON(code, httpresponse.Error{
			Code:    code,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Reservation confirmed successfully",
		"data":    res,
	})
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: user ID not found in context",
		})
	}

	res, err := h.service.GetMyReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Could not retrieve reservations",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "My reservations retrieved successfully",
		"data":    res,
	})
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	res, err := h.service.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Could not retrieve reservations",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "All reservations retrieved successfully",
		"data":    res,
	})
}

func (h *ReservationHandler) CancelReservation(c echo.Context) error {
	userIDVal := c.Get("user_id")
	userID, ok := userIDVal.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: user ID not found in context",
		})
	}

	roleVal := c.Get("role")
	role, ok := roleVal.(string)
	if !ok {
		return c.JSON(http.StatusUnauthorized, httpresponse.Error{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized: role not found in context",
		})
	}

	idStr := c.Param("id")
	reservationID, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid reservation ID",
		})
	}

	err = h.service.CancelReservation(userID, role, uint(reservationID))
	if err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, ErrReservationNotFound) {
			code = http.StatusNotFound
		} else if errors.Is(err, ErrForbidden) {
			code = http.StatusForbidden
		} else if errors.Is(err, ErrReservationNotActive) {
			code = http.StatusBadRequest
		}

		return c.JSON(code, httpresponse.Error{
			Code:    code,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Reservation cancelled successfully",
	})
}
