package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBTimezone string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file")
	}

	config := &Config{
		ServerPort: os.Getenv("SERVER_PORT"),

		DBHost:     os.Getenv("DB_HOST"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBPort:     os.Getenv("DB_PORT"),
		DBTimezone: os.Getenv("DB_TIMEZONE"),
	}

	return config, nil
}
