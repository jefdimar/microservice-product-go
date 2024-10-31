package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostgresProduct struct {
	ID          uint      `json:"id" gorm:"primary_key" example:"1"`
	Name        string    `json:"name" example:"Product Name"`
	Description string    `json:"description" example:"Product Description"`
	Price       float64   `json:"price" example:"100.00"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type MongoProduct struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"6123456789abcdef0123456"`
	Name        string             `json:"name" example:"Product Name"`
	Description string             `json:"description" example:"Product Description"`
	Price       float64            `json:"price" example:"100.00"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type Product struct {
	ID          interface{}	`json:"id"`
	Name        string    	`json:"name"`
	Description string    	`json:"description"`
	Price       float64   	`json:"price"`
	CreatedAt   time.Time 	`json:"created_at"`
	UpdatedAt   time.Time 	`json:"updated_at"`
}

func (p PostgresProduct) ToCommon() Product {
	return Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (p MongoProduct) ToCommon() Product {
	return Product{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}