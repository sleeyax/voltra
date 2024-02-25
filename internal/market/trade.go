package market

type BuyOrder struct {
	OrderID          int64
	Symbol           string
	TransactionTime  int64
	Price            float64
	ExecutedQuantity string
}

type BuyOrderMap map[string]BuyOrder
