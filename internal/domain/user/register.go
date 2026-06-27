package user

import (
	"spotsync/internal/auth"
	"spotsync/internal/middleware"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// RegisterRoutes initializes user dependencies (Repository, Service, Handler) and registers routes with Echo.
func RegisterRoutes(e *echo.Echo, db *gorm.DB, jwtService auth.JWTService) {
	// 1. Initialize layers
	userRepo := NewUserRepository(db)
	userService := NewUserService(userRepo, jwtService)
	userHandler := NewUserHandler(userService)

	// 2. Auth Middleware
	authMiddleware := middleware.AuthMiddleware(jwtService)

	// 3. Register routes
	g := e.Group("/api/v1/auth")
	g.POST("/register", userHandler.Register)
	g.POST("/login", userHandler.Login)
	g.GET("/me", userHandler.GetMe, authMiddleware)
	g.POST("/refresh", userHandler.RefreshToken, authMiddleware)
}
