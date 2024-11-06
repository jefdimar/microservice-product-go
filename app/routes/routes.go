package routes

import (
	"go-microservice-product-porto/app/controllers"
	"go-microservice-product-porto/app/repositories"
	"go-microservice-product-porto/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine) {
	// Default route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})

	// Swagger documentation
	r.GET("/doc", func(c *gin.Context) {
		c.Request.URL.Path = "/doc/index.html"
		r.HandleContext(c)
	})
	r.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize repository
	productRepo := repositories.NewProductRepository()

	// Initialize business
	productUsecase := usecase.NewProductUsecase(productRepo)

	// Initialize controller
	productController := controllers.NewProductController(productUsecase)

	// API Routes
	api := r.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("/", productController.Create)
			products.GET("/", productController.GetAll)
			products.GET("/:id", productController.GetByID)
			products.PUT("/:id", productController.Update)
			products.DELETE("/:id", productController.Delete)
		}
	}
}
