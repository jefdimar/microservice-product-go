package config

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func InitMongoDB(cfg *Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	// Construct MongoDB connection URI
	mongoURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=admin",
		cfg.MongoUser,
		cfg.MongoPassword,
		cfg.MongoHost,
		cfg.MongoPort,
		cfg.MongoDBName,
	)

	// Create client options
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Ping the database
	if err = client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	// Set the database instance
	MongoDB = client.Database(cfg.MongoDBName)

	return nil
}
