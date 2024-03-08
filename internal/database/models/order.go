package models

import (
	"github.com/sleeyax/voltra/internal/market"
	"gorm.io/gorm"
)

type OrderType string

const (
	BuyOrder  OrderType = "buy"
	SellOrder OrderType = "sell"
)

type Order struct {
	gorm.Model
	market.Order

	// Required field to indicate which market the order is for.
	Market string

	// Required field to indicate the type of order.
	Type OrderType

	// Required field to indicate the volume of the symbol.
	Volume float64

	// Optional field to store the take profit.
	// This field may be updated when trailing stop loss is used.
	// This field is only set when the type is a buy order.
	TakeProfit *float64

	// Optional field to store the stop loss.
	// This field may be updated when trailing stop loss is used.
	// This field is only set when the type is a buy order.
	StopLoss *float64

	// Optional field for the estimated profit.
	// This field is only set when the type is a sell order.
	PriceChangePercentage *float64

	// Optional field for the estimated profit or loss.
	// This field is only set when the type is a sell order.
	EstimatedProfitLoss *float64

	// Whether the order is a dummy/fake order, created in test mode.
	IsTestMode bool
}
