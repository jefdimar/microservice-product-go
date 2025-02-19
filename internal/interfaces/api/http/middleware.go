package http

import (
	"time"

	"go-microservice-product-porto/pkg/logger"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add authentication logic here
		// Example: JWT validation
		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		logger.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Str("ip", c.ClientIP()).
			Str("user-agent", c.Request.UserAgent()).
			Msg("incoming request")

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logEvent := logger.Info()
		if statusCode >= 400 {
			logEvent = logger.Error()
		}

		logEvent.
			Int("status", statusCode).
			Str("method", c.Request.Method).
			Str("path", path).
			Str("query", query).
			Str("ip", c.ClientIP()).
			Str("latency", latency.String()).
			Int("size", c.Writer.Size()).
			Msg("request completed")
	}
}
