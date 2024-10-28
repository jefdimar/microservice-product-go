package main

import (
	"go-microservice-product-porto/app/routes"
	"go-microservice-product-porto/config"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize database
	err = config.InitDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Set up Gin
	r := gin.Default()

	// Initialize routes
	routes.SetupRoutes(r)

	// Run server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running on %s", cfg.ServerPort)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
