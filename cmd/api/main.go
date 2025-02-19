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
	"go-microservice-product-porto/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init("debug")
	logger.Info().Msg("Initializing application...")

	// Load configuration
	logger.Info().Msg("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to load configuration")

		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize MongoDB client
	logger.Info().Msg("Initializing MongoDB client...")
	mongoClient, err := mongodb.InitMongoDB(mongodb.MongoConfig{
		Host:     cfg.MongoHost,
		Port:     cfg.MongoPort,
		User:     cfg.MongoUser,
		Password: cfg.MongoPassword,
		DBName:   cfg.MongoDBName,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to initialize MongoDB client")

		log.Fatal("Failed to initialize MongoDB client:", err)
	}

	// Initialize repository
	logger.Info().Msg("Initializing repository...")
	productRepo := mongodb.NewProductRepository(mongoClient)

	// Initialize Redis cache
	logger.Info().Msg("Initializing Redis cache...")
	cacheService, err := cache.NewCacheService(redis.RedisConfig{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to initialize Redis cache")

		log.Fatal("Failed to initialize Redis cache:", err)
	}

	// Initialize event handler
	logger.Info().Msg("Initializing event handler...")
	eventHandler := eventhandlers.NewProductEventHandler(cacheService, productRepo)

	// Initialize command handler
	logger.Info().Msg("Initializing command handler...")
	commandHandler := commands.NewProductCommandHandler(productRepo, eventHandler, cacheService)

	// Initialize query handler
	logger.Info().Msg("Initializing query handler...")
	queryHandler := queries.NewProductQueryHandler(productRepo, cacheService)

	// Initialize HTTP handler
	logger.Info().Msg("Initializing HTTP handler...")
	productHandler := http.NewProductHandler(commandHandler, queryHandler)

	// Setup router
	logger.Info().Msg("Setting up router...")
	router := http.SetupRouter(productHandler)

	// Start server
	logger.Info().Msg("Starting server...")
	if err := router.Run(cfg.ServerAddress); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
