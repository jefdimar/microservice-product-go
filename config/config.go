package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort 			string
	TrustedProxies 	string

	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	PostgresPort     string
	PostgresTimezone string

	MongoHost 		string
	MongoPort 		string
	MongoUser 		string
	MongoPassword string
	MongoDBName 	string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	config := &Config{
		ServerPort:     os.Getenv("SERVER_PORT"),
		TrustedProxies: os.Getenv("TRUSTED_PROXIES"),
		PostgresHost:         os.Getenv("POSTGRES_HOST"),
		PostgresUser:         os.Getenv("POSTGRES_USER"),
		PostgresPassword:         os.Getenv("POSTGRES_PASSWORD"),
		PostgresDBName:     os.Getenv("POSTGRES_DBNAME"),
		PostgresPort: os.Getenv("POSTGRES_PORT"),
		MongoHost: os.Getenv("MONGO_HOST"),
		MongoPort: os.Getenv("MONGO_PORT"),
		MongoUser: os.Getenv("MONGO_USER"),
		MongoPassword: os.Getenv("MONGO_PASSWORD"),
		MongoDBName: os.Getenv("MONGO_DB_NAME"),
	}

	return config, nil
}