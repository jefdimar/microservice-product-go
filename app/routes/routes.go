package routes

import (
	"go-microservice-product-porto/app/controllers"
	"go-microservice-product-porto/app/handlers"
	"go-microservice-product-porto/app/repositories"
	"go-microservice-product-porto/app/services"
	"go-microservice-product-porto/app/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(r *gin.Engine, redisClient *redis.Client) {
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

	// Initialize all components
	responseHandler := handlers.NewResponseHandler()
	cacheService := services.NewCacheService(redisClient)
	productRepo := repositories.NewProductRepository(cacheService)
	productUsecase := usecase.NewProductUsecase(productRepo, cacheService)
	productController := controllers.NewProductController(productUsecase, responseHandler)

	api := r.Group("/api")
	{
		products := api.Group("/products")
		{
			products.POST("/", productController.Create)
			products.GET("/", productController.GetAll)
			products.GET("/:id", productController.GetByID)
			products.PATCH("/:id", productController.Update)
			products.DELETE("/:id", productController.Delete)
		}
	}
}
