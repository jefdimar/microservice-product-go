package models

import "time"

type Product struct {
	ID          uint      `json:"id" gorm:"primary_key" example:"1"`
	Name        string    `json:"name" example:"Product Name"`
	Description string    `json:"description" example:"Product Description"`
	Price       float64   `json:"price" example:"100.00"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
