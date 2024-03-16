package database

import (
	"github.com/sleeyax/voltra/internal/config"
	"github.com/sleeyax/voltra/internal/database/models"
	"github.com/sleeyax/voltra/internal/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
	"time"
)

type SqliteDatabase struct {
	db *gorm.DB
}

var _ Database = (*SqliteDatabase)(nil)

func NewSqliteDatabase(fileName string, options config.LoggingOptions) *SqliteDatabase {
	var logLevel config.LogLevel
	if !options.Enable || options.EnableStructuredLogging {
		logLevel = config.SilentLevel
	} else if options.DatabaseLogLevel != "" {
		logLevel = options.DatabaseLogLevel
	} else {
		logLevel = options.LogLevel
	}

	customLoggerConfig := logger.Config{
		SlowThreshold:             time.Second * 1,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      options.EnableStructuredLogging,
		Colorful:                  !options.EnableStructuredLogging,
		LogLevel:                  toGORMLogLevel(logLevel),
	}

	_ = storage.CreateDataDirectoryIfNotExists()

	db, err := gorm.Open(sqlite.Open(filepath.Join(storage.DataPath, fileName)), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			customLoggerConfig,
		),
	})
	if err != nil {
		panic("failed to connect to the local database: " + err.Error())
	}

	_ = db.AutoMigrate(&models.Order{})
	_ = db.AutoMigrate(&models.Cache{})

	return &SqliteDatabase{db: db}
}

// Converts a config.LogLevel to a gorm logger.LogLevel.
func toGORMLogLevel(level config.LogLevel) logger.LogLevel {
	switch level {
	case config.WarnLevel:
		return logger.Warn
	case config.ErrorLevel:
		return logger.Error
	case config.SilentLevel:
		return logger.Silent
	case config.DebugLevel:
		fallthrough
	case config.InfoLevel:
		fallthrough
	default:
		return logger.Info
	}
}

func (d *SqliteDatabase) SaveOrder(order models.Order) {
	d.db.Save(&order)
}

func (d *SqliteDatabase) HasOrder(orderType models.OrderType, market, symbol string) bool {
	var count int64
	d.db.Model(&models.Order{}).Where("type = ? AND market = ? AND symbol = ?", orderType, market, symbol).Count(&count)
	return count > 0
}

func (d *SqliteDatabase) CountOrders(orderType models.OrderType, market string) int64 {
	var count int64
	d.db.Model(&models.Order{}).Where("type = ? AND market = ?", orderType, market).Count(&count)
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
