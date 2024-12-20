package main

import (
	"fmt"
	"go-microservice-product-porto/app/routes"
	"go-microservice-product-porto/config"
	_ "go-microservice-product-porto/docs"
	"log"

	"github.com/gin-gonic/gin"
)

// @title           Product API
// @version         1.0
// @description     A Product microservice API in Go using Gin framework.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize Redis using config
	redisClient := config.InitRedis(cfg)

	// Initialize database
	err = config.InitDatabases(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Set up Gin
	r := gin.Default()

	// Initialize routes with Redis client
	routes.SetupRoutes(r, redisClient)

	// Run server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server is running on %s", cfg.ServerPort)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
