package helpers

import (
	"fmt"
	"strings"
)

func FormatPrice(price float64) string {
	// Convert price to string with 2 decimal places
	priceStr := fmt.Sprintf("%.2f", price)

	// Split the integer and decimal parts
	parts := strings.Split(priceStr, ".")
	intPart := parts[0]

	// Add thousand separators
	var result []byte
	for i, j := len(intPart)-1, 0; i >= 0; i, j = i-1, j+1 {
		if j > 0 && j%3 == 0 {
			result = append([]byte{'.'}, result...)
		}
		result = append([]byte{intPart[i]}, result...)
	}

	// Combine with currency code and decimal part
	return fmt.Sprintf("IDR %s,%s", string(result), parts[1])
}
