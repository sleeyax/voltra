package database

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type SqliteDatabase struct {
	db *gorm.DB
}

var _ Database = (*SqliteDatabase)(nil)

func NewSqliteDatabase(dsn string, options config.LoggingOptions) *SqliteDatabase {
	customLoggerConfig := logger.Config{
		SlowThreshold:             time.Second * 1,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      options.EnableStructuredLogging,
		Colorful:                  !options.EnableStructuredLogging,
	}

	if !options.Enable || options.EnableStructuredLogging {
		customLoggerConfig.LogLevel = logger.Silent
	}

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			customLoggerConfig,
		),
	})
	if err != nil {
		panic("failed to connect to the local database")
	}

	_ = db.AutoMigrate(&models.Order{})
	_ = db.AutoMigrate(&models.Cache{})

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

func (d *SqliteDatabase) SaveCache(cache models.Cache) {
	d.db.Save(&cache)
}

func (d *SqliteDatabase) GetCache(symbol string) (models.Cache, bool) {
	var cache models.Cache
	if err := d.db.Where("symbol = ?", symbol).First(&cache).Error; err != nil {
		return cache, false
	}
	return cache, true
}

func (d *SqliteDatabase) GetLastOrder(orderType models.OrderType, market, symbol string) (models.Order, bool) {
	var order models.Order
	if err := d.db.Where("type = ? AND market = ? AND symbol = ?", orderType, market, symbol).Last(&order).Error; err != nil {
		return order, false
	}
	return order, true
}
