package parking_zone

import (
	"net/http"
	"strconv"

	"spotsync/internal/domain/parking_zone/dto"
	"spotsync/internal/httpResponse"

	"github.com/labstack/echo/v4"
)

type ParkingZoneHandler struct {
	service ParkingZoneService
}

func NewParkingZoneHandler(service ParkingZoneService) *ParkingZoneHandler {
	return &ParkingZoneHandler{service: service}
}

func (h *ParkingZoneHandler) CreateZone(c echo.Context) error {
	var req dto.CreateParkingZoneRequest
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

	res, err := h.service.CreateZone(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Parking zone created successfully",
		"data":    res,
	})
}

func (h *ParkingZoneHandler) GetAllZones(c echo.Context) error {
	res, err := h.service.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, httpresponse.Error{
			Code:    http.StatusInternalServerError,
			Message: "Could not retrieve parking zones",
			Details: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zones retrieved successfully",
		"data":    res,
	})
}

func (h *ParkingZoneHandler) GetZoneByID(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, httpresponse.Error{
			Code:    http.StatusBadRequest,
			Message: "Invalid parking zone ID",
		})
	}

	res, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, httpresponse.Error{
			Code:    http.StatusNotFound,
			Message: "Parking zone not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    res,
	})
}
