package product

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Product, error)
	FindAll(ctx context.Context) ([]*Product, error)
	Update(ctx context.Context, product *Product) error
	Create(ctx context.Context, product *Product) error
	Delete(ctx context.Context, id string) error
	Search(ctx context.Context, name string, minPrice, maxPrice float64) ([]*Product, error)
}
