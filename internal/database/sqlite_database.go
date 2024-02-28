package database

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDatabase struct {
	db *gorm.DB
}

var _ Database = (*SqliteDatabase)(nil)

func NewSqliteDatabase(dsn string) *SqliteDatabase {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the local database")
	}

	_ = db.AutoMigrate(&models.Order{})

	return &SqliteDatabase{db: db}
}

func (d *SqliteDatabase) SaveOrder(order models.Order) {
	d.db.Save(&order)
}

func (d *SqliteDatabase) HasOrder(orderType models.OrderType, market, symbol string) bool {
	var count int64
	d.db.Where("type = ? AND market = ? AND symbol = ?", orderType, market, symbol).Count(&count)
	return count > 0
}

func (d *SqliteDatabase) CountOrders(orderType models.OrderType, market string) int64 {
	var count int64
	d.db.Where("type = ? AND market = ?", orderType, market).Count(&count)
	return count
}

func (d *SqliteDatabase) GetOrders(orderType models.OrderType, market string) []models.Order {
	var orders []models.Order
	d.db.Where("type = ? AND market = ?", orderType, market).Find(&orders)
	return orders
}

func (d *SqliteDatabase) DeleteOrder(order models.Order) {
	d.db.Delete(&order)
}
