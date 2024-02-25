package models

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"gorm.io/gorm"
)

type OrderType string

const (
	BuyOrder  OrderType = "buy"
	SellOrder OrderType = "sell"
)

type Order struct {
	gorm.Model
	market.BuyOrder
	Market     string
	Type       OrderType
	Volume     float64
	IsTestMode bool
}
