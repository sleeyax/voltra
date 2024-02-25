package market

import "time"

type BuyOrder struct {
	OrderID         int64
	Symbol          string
	TransactionTime time.Time
	Price           float64
}

type BuyOrderMap map[string]BuyOrder
