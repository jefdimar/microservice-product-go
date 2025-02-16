package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *ProductHandler) *gin.Engine {
	router := gin.Default()

	// API routes
	v1 := router.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			// Add this new route
			products.GET("/search", handler.SearchProducts)

			// Existing routes remain unchanged
			products.POST("/", handler.CreateProduct)
			products.GET("/", handler.ListProducts)
			products.GET("/:id", handler.GetProduct)
			products.PATCH("/:id/stock", handler.UpdateStock)
			products.DELETE("/:id", handler.DeleteProduct)
		}
	}

	return router
}
