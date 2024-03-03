package database

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database/models"
)

type Database interface {
	SaveOrder(order models.Order)
	HasOrder(orderType models.OrderType, market, symbol string) bool
	CountOrders(orderType models.OrderType, market string) int64
	GetOrders(orderType models.OrderType, market string) []models.Order
	GetLastOrder(orderType models.OrderType, market, symbol string) (models.Order, bool)
	DeleteOrder(order models.Order)
	SaveCache(cache models.Cache)
	GetCache(symbol string) (models.Cache, bool)
}
