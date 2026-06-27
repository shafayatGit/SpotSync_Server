package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DSN       string
	JWTSecret string
}

// LoadConfig loads the configuration from .env file or environment variables.
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it, reading from environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN environment variable is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}

	return &Config{
		Port:      port,
		DSN:       dsn,
		JWTSecret: jwtSecret,
	}
}
