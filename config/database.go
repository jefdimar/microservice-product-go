package config

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DBConnections struct {
	MongoDB   *mongo.Database
	PostgreDB *gorm.DB
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

	return nil
}
