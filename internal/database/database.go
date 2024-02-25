package database

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	SaveOrder(order models.Order)
	HasOrder(orderType models.OrderType, market, symbol string) bool
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
	d.db.Create(&order)
}

func (d *LocalDatabase) HasOrder(orderType models.OrderType, market, symbol string) bool {
	var count int64
	d.db.Where("type = ? AND market = ? AND symbol = ?", orderType, market, symbol).Count(&count)
	return count > 0
}
