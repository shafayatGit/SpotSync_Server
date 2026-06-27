package parking_zone

import (
	"net/http"
	"strconv"

	"spotsync/internal/domain/parking_zone/dto"

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

	res, err := h.service.CreateZone(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": err.Error(),
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
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Could not retrieve parking zones",
			"errors":  err.Error(),
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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid parking zone ID",
		})
	}

	res, err := h.service.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success": false,
			"message": "Parking zone not found",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    res,
	})
}
