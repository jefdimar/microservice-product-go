package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	skuPrefix = "PRD-"
	skuLength = 8
)

func GenerateSKU() string {
	rand.Seed(time.Now().UnixNano())

	numbers := make([]byte, skuLength)
	for i := 0; i < skuLength; i++ {
		numbers[i] = byte(rand.Intn(10) + '0')
	}

	return fmt.Sprintf("%s-%s", skuPrefix, string(numbers))
}
