package reservation

import (
	"spotsync/internal/auth"
	"spotsync/internal/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func ReservationRoute(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	reservationRepo := NewReservationRepository(db)
	reservationService := NewReservationService(reservationRepo)
	reservationHandler := NewReservationHandler(reservationService)

	authMiddleware := middleware.AuthMiddleware(jwtService)
	adminMiddleware := middleware.RoleMiddleware("admin")

	g := e.Group("/api/v1/reservations")
	g.POST("", reservationHandler.ReserveSpot, authMiddleware)
	g.GET("/my-reservations", reservationHandler.GetMyReservations, authMiddleware)
	g.DELETE("/:id", reservationHandler.CancelReservation, authMiddleware)
	g.GET("", reservationHandler.GetAllReservations, authMiddleware, adminMiddleware)
}
