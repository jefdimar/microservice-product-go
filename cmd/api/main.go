package main

import (
	"log"

	"go-microservice-product-porto/internal/application/commands"
	eventhandlers "go-microservice-product-porto/internal/application/event_handlers"
	"go-microservice-product-porto/internal/application/queries"
	"go-microservice-product-porto/internal/infrastructure/cache"
	"go-microservice-product-porto/internal/infrastructure/persistence/mongodb"
	"go-microservice-product-porto/internal/infrastructure/persistence/redis"
	"go-microservice-product-porto/internal/interfaces/api/http"

	"go-microservice-product-porto/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize MongoDB client
	mongoClient, err := mongodb.InitMongoDB(mongodb.MongoConfig{
		Host:     cfg.MongoHost,
		Port:     cfg.MongoPort,
		User:     cfg.MongoUser,
		Password: cfg.MongoPassword,
		DBName:   cfg.MongoDBName,
	})
	if err != nil {
		log.Fatal("Failed to initialize MongoDB client:", err)
	}

	// Initialize repository
	productRepo := mongodb.NewProductRepository(mongoClient)

	// Initialize Redis cache
	cacheService, err := cache.NewCacheService(redis.RedisConfig{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})
	if err != nil {
		log.Fatal("Failed to initialize Redis cache:", err)
	}

	// Initialize event handler
	eventHandler := eventhandlers.NewProductEventHandler(cacheService, productRepo)

	// Initialize command handler
	commandHandler := commands.NewProductCommandHandler(productRepo, eventHandler, cacheService)

	// Initialize query handler
	queryHandler := queries.NewProductQueryHandler(productRepo, cacheService)

	// Initialize HTTP handler
	productHandler := http.NewProductHandler(commandHandler, queryHandler)

	// Setup router
	router := http.SetupRouter(productHandler)

	// Start server
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
