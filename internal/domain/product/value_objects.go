package product

type Money struct {
	Amount   float64
	Currency string
}

func NewMoney(amount float64, currency string) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

type Stock struct {
	Quantity int
	Unit     string
}

func NewStock(quantity int, unit string) Stock {
	return Stock{
		Quantity: quantity,
		Unit:     unit,
	}
}
