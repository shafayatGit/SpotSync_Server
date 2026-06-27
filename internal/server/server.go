package server

import (
	"log"
	"net/http"
	"time"

	"spotsync/internal/auth"
	"spotsync/internal/config"
	"spotsync/internal/domain/parking_zone"
	"spotsync/internal/domain/reservation"
	"spotsync/internal/domain/user"

	"github.com/go-playground/validator/v10"
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

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// NewServer initializes the database connection, GORM configuration, and Echo server.
func NewServer(cfg *config.Config) *Server {
	// 1. Initialize database connection
	log.Printf("Connecting to database...")
	// We use PreferSimpleProtocol: true to prevent Neon connection pooler (PgBouncer) issues with GORM migrations.
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  cfg.DSN,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
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

	// Run migrations
	log.Println("Running database migrations...")
	if err := db.AutoMigrate(&user.User{}, &parking_zone.ParkingZone{}, &reservation.Reservation{}); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
	log.Println("Database migrations completed successfully.")

	// 2. Initialize Echo
	e := echo.New()

	// Setup custom validation
	e.Validator = &CustomValidator{validator: validator.New()}

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

	// 4. Initialize Core Services
	jwtService := auth.NewJWTService(cfg.JWTSecret)

	// 5. Register Routes (Module-level Dependency Injection)
	user.RegisterRoutes(e, db, jwtService)
	parking_zone.ParkingZoneRoute(e, db, jwtService)

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
