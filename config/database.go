package config

import (
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DBConnections struct {
	MongoDB   *mongo.Database
	PostgreDB *gorm.DB
	Redis     *redis.Client
}

var DBConn DBConnections

func InitDatabases(cfg *Config) error {
	// Initialize MongoDB
	mongoErr := InitMongoDB(cfg)
	if mongoErr != nil {
		return mongoErr
	}
	DBConn.MongoDB = MongoDB

	// Initialize PostgreSQL
	postgresDB, postgresErr := InitPostgres(cfg)
	if postgresErr != nil {
		return postgresErr
	}
	DBConn.PostgreDB = postgresDB

	// Initialize Redis
	redisClient := InitRedis(cfg)
	DBConn.Redis = redisClient

	return nil
}
