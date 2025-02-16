package product

import "errors"

var (
	ErrProductNotFound      = errors.New("product not found")
	ErrInvalidProduct       = errors.New("invalid product")
	ErrInvalidStock         = errors.New("invalid stock value")
	ErrProductAlreadyExists = errors.New("product already exists")
)
