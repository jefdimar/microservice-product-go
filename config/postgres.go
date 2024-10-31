package config

import (
	"fmt"
	"go-microservice-product-porto/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres(cfg *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
        cfg.PostgresHost,
        cfg.PostgresUser,
        cfg.PostgresPassword,
        cfg.PostgresDBName,
        cfg.PostgresPort,
    )

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PostgreSQL: %v", err)
    }

    // Auto Migrate your models here
    db.AutoMigrate(&models.PostgresProduct{})

    return db, nil
}
