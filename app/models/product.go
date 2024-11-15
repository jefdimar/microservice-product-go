package models

import (
	"go-microservice-product-porto/app/helpers"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostgresProduct struct {
	ID          uint      `json:"id" gorm:"primary_key" example:"1"`
	SKU         string    `json:"sku" gorm:"unique" example:"PRD-12345678"`
	Name        string    `json:"name" example:"Product Name"`
	Description string    `json:"description" example:"Product Description"`
	Price       float64   `json:"price" example:"100.00"`
	Stock       int       `json:"stock" example:"100"`
	IsActive    bool      `json:"is_active" example:"true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MongoProduct struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"6123456789abcdef0123456"`
	SKU         string             `json:"sku" bson:"sku" example:"PRD-12345678"`
	Name        string             `json:"name" example:"Product Name"`
	Description string             `json:"description" example:"Product Description"`
	Price       float64            `json:"price" example:"100.00"`
	Stock       int                `json:"stock" example:"100"`
	IsActive    bool               `json:"is_active" example:"true"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Product struct {
	ID                 interface{} `json:"id"`
	SKU                string      `json:"sku"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	Price              float64     `json:"price"`
	FormattedPrice     string      `json:"formatted_price"`
	Stock              int         `json:"stock"`
	IsActive           bool        `json:"is_active"`
	CreatedAt          time.Time   `json:"-"` // Hide the original field
	UpdatedAt          time.Time   `json:"-"` // Hide the original field
	FormattedCreatedAt string      `json:"created_at"`
	FormattedUpdatedAt string      `json:"updated_at"`
}

type ProductUpdate struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

func (p PostgresProduct) ToCommon() Product {
	return Product{

		ID:                 p.ID,
		SKU:                p.SKU,
		Name:               p.Name,
		Description:        p.Description,
		Price:              p.Price,
		Stock:              p.Stock,
		IsActive:           p.IsActive,
		FormattedCreatedAt: helpers.FormatDateTime(p.CreatedAt),
		FormattedUpdatedAt: helpers.FormatDateTime(p.UpdatedAt),
	}
}

func (p MongoProduct) ToCommon() Product {
	return Product{

		ID:                 p.ID,
		SKU:                p.SKU,
		Name:               p.Name,
		Description:        p.Description,
		Price:              p.Price,
		Stock:              p.Stock,
		IsActive:           p.IsActive,
		FormattedCreatedAt: helpers.FormatDateTime(p.CreatedAt),
		FormattedUpdatedAt: helpers.FormatDateTime(p.UpdatedAt),
	}
}
