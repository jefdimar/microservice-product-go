package routes

import (
	"go-microservice-product-porto/app/business"
	"go-microservice-product-porto/app/controllers"
	"go-microservice-product-porto/app/repositories"
	"go-microservice-product-porto/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Default route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Product Service"})
	})

	// Initialize repository
	productRepo := repositories.NewProductRepository(config.DB)

	// Initialize business
	productBusiness := business.NewProductBusiness(productRepo)

	// Initialize controller
	productController := controllers.NewProductController(productBusiness)

	// Group routes
	api := r.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("/", productController.Create)
			products.GET("/", productController.GetAll)
			products.GET("/:id", productController.GetByID)
		}
	}
}
