package validation

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateObjectID(id string) error {
	if !primitive.IsValidObjectID(id) {
		return fmt.Errorf("invalid id format")
	}
	return nil
}
