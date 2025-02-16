package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Price       float64            `bson:"price" json:"price"`
	Stock       int                `bson:"stock" json:"stock"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

func NewProduct(name, description string, price float64, stock int) *Product {
	return &Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func (p *Product) UpdateStock(newStock int) error {
	if newStock < 0 {
		return ErrInvalidStock
	}
	p.Stock = newStock
	p.UpdatedAt = time.Now()
	return nil
}

func (p *Product) IsValid() bool {
	return p.Name != "" && p.Price > 0 && p.Stock >= 0
}
