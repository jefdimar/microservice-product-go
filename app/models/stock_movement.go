package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockMovement struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID     primitive.ObjectID `bson:"product_id" json:"product_id"`
	Type          string             `bson:"type" json:"type"`
	Quantity      int                `bson:"quantity" json:"quantity"`
	PreviousStock int                `bson:"previous_stock" json:"previous_stock"`
	NewStock      int                `bson:"new_stock" json:"new_stock"`
	Reason        string             `bson:"reason" json:"reason"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}
