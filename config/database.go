package config

import (
	"fmt"
	"go-microservice-product-porto/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase(config *Config) error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
	config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBTimezone)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v", err)
	}

	// Auto migrate models
	err = DB.AutoMigrate(&models.Product{})
	if err != nil {
		return fmt.Errorf("Failed to migrate models: %v", err)
	}
	return nil
}
