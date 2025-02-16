package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func InitMongoDB(cfg MongoConfig) (*mongo.Client, error) {
	var uri string

	// Build connection string based on credentials presence
	if cfg.User != "" && cfg.Password != "" {
		uri = fmt.Sprintf("mongodb://%s:%s@%s:%s", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	} else {
		uri = fmt.Sprintf("mongodb://%s:%s", cfg.Host, cfg.Port)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri).SetDirect(true)

	// Add authentication if credentials are provided
	if cfg.User != "" && cfg.Password != "" {
		credential := options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		}
		clientOptions.SetAuth(credential)
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
