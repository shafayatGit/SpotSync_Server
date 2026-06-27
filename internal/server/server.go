package server

import (
	"log"
	"net/http"
	"time"

	"spotsync/internal/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	Echo *echo.Echo
	DB   *gorm.DB
	Cfg  *config.Config
}

// NewServer initializes the database connection, GORM configuration, and Echo server.
func NewServer(cfg *config.Config) *Server {
	// 1. Initialize database connection
	log.Printf("Connecting to database...")
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connection established successfully.")

	// Configure connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance from GORM: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 2. Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 3. Default / Health check route
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"message": "Welcome to SpotSync API",
			"data": map[string]string{
				"status": "healthy",
				"time":   time.Now().Format(time.RFC3339),
			},
		})
	})

	return &Server{
		Echo: e,
		DB:   db,
		Cfg:  cfg,
	}
}

// Start starts the HTTP server listening on the configured port.
func (s *Server) Start() {
	serverAddr := ":" + s.Cfg.Port
	log.Printf("Starting server on %s", serverAddr)
	if err := s.Echo.Start(serverAddr); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
