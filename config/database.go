package config

import (
	"go-microservice-product-porto/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitConfig() {
	var err error
	dsn := "host=localhost user=postgres password=postgres dbname=product_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto migrate models
	DB.AutoMigrate(&models.Product{})
}
