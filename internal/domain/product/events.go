package product

type Event interface {
	GetEventType() string
}

type ProductCreatedEvent struct {
	Product *Product
}

func (e ProductCreatedEvent) GetEventType() string {
	return "product.created"
}

type ProductStockUpdatedEvent struct {
	Product  *Product
	OldStock int
	NewStock int
}

func (e ProductStockUpdatedEvent) GetEventType() string {
	return "product.stock.updated"
}

type ProductDeletedEvent struct {
	ProductID string
}

func (e ProductDeletedEvent) GetEventType() string {
	return "product.deleted"
}
