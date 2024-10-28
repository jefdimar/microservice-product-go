package main

import (
	"go-microservice-product-porto/app/routes"
	"go-microservice-product-porto/config"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize config
	config.InitConfig()

	// Set up Gin
	r := gin.Default()

	// Initialize routes
	routes.SetupRoutes(r)

	// Run server
	r.Run(":8080")
}
