package product

import "context"

type Repository interface {
	Create(context.Context, *Product) error
	FindByID(context.Context, string) (*Product, error)
	FindAll(ctx context.Context, page, pageSize int, sortBy, sortDir string) ([]*Product, int64, error)
	Update(context.Context, *Product) error
	Delete(context.Context, string) error
	Search(context.Context, string, float64, float64) ([]*Product, error)
}
