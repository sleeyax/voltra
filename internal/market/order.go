package market

import "time"

type Order struct {
	OrderID         int64
	Symbol          string
	TransactionTime time.Time
	Price           float64
}
