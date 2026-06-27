package parking_zone

import (
	"spotsync/internal/auth"
	"spotsync/internal/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ParkingZoneRoute(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	parkingZoneRepo := NewParkingZoneRepository(db)
	parkingZoneService := NewParkingZoneService(parkingZoneRepo)
	parkingZoneHandler := NewParkingZoneHandler(parkingZoneService)

	authMiddleware := middleware.AuthMiddleware(jwtService)
	adminMiddleware := middleware.RoleMiddleware("admin")

	g := e.Group("/api/v1/zones")
	g.POST("", parkingZoneHandler.CreateZone, authMiddleware, adminMiddleware)
	g.GET("", parkingZoneHandler.GetAllZones)
	g.GET("/:id", parkingZoneHandler.GetZoneByID)
}
