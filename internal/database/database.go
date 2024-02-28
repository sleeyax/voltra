package database

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	SaveOrder(order models.Order)
	HasOrder(orderType models.OrderType, market, symbol string) bool
	CountOrders(orderType models.OrderType, market string) int64
	GetOrders(orderType models.OrderType, market string) []models.Order
	DeleteOrder(order models.Order)
}

const LocalDatabasePath = "data.db"

type LocalDatabase struct {
	db *gorm.DB
}

var _ Database = (*LocalDatabase)(nil)

func NewLocalDatabase() *LocalDatabase {
	db, err := gorm.Open(sqlite.Open(LocalDatabasePath), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the local database")
	}

	db.AutoMigrate(&models.Order{})

	return &LocalDatabase{db: db}
}

func (d *LocalDatabase) SaveOrder(order models.Order) {
	d.db.Save(&order)
}

func (d *LocalDatabase) HasOrder(orderType models.OrderType, market, symbol string) bool {
	var count int64
	d.db.Where("type = ? AND market = ? AND symbol = ?", orderType, market, symbol).Count(&count)
	return count > 0
}

func (d *LocalDatabase) CountOrders(orderType models.OrderType, market string) int64 {
	var count int64
	d.db.Where("type = ? AND market = ?", orderType, market).Count(&count)
	return count
}

func (d *LocalDatabase) GetOrders(orderType models.OrderType, market string) []models.Order {
	var orders []models.Order
	d.db.Where("type = ? AND market = ?", orderType, market).Find(&orders)
	return orders
}

func (d *LocalDatabase) DeleteOrder(order models.Order) {
	d.db.Delete(&order)
}
