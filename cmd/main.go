package main

import (
	"spotsync/internal/config"
	"spotsync/internal/server"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize and start server
	srv := server.NewServer(cfg)
	srv.Start()
}
