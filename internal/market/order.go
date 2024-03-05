package market

import "time"

type Order struct {
	OrderID         int64
	TransactionTime time.Time
	Symbol          string
	Price           float64
}
