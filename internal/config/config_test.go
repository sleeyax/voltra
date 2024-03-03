package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoad(t *testing.T) {
	config, err := Load("../../config.example.yml")

	assert.Nil(t, err)

	assert.Equal(t, true, config.EnableTestMode)

	assert.Equal(t, true, config.LoggingOptions.Enable)
	assert.Equal(t, false, config.LoggingOptions.EnableStructuredLogging)

	assert.Equal(t, "PASTE_YOUR_ACCESS_KEY_HERE", config.Markets.Binance.AccessKey)
	assert.Equal(t, "PASTE_YOUR_SECRET_KEY_HERE", config.Markets.Binance.SecretKey)

	assert.Equal(t, "USDT", config.TradingOptions.PairWith)
	assert.Equal(t, float64(15), config.TradingOptions.Quantity)
	assert.Equal(t, 3, config.TradingOptions.MaxCoins)
	assert.Equal(t, 2, config.TradingOptions.TimeDifference)
	assert.Equal(t, 10, config.TradingOptions.RecheckInterval)
	assert.Equal(t, 10, config.TradingOptions.SellTimeout)
	assert.Equal(t, float64(10), config.TradingOptions.ChangeInPrice)
	assert.Equal(t, float64(5), config.TradingOptions.StopLoss)
	assert.Equal(t, 0.8, config.TradingOptions.TakeProfit)
	assert.Equal(t, 0.075, config.TradingOptions.TradingFee)
	assert.Equal(t, 0, config.TradingOptions.CoolOffDelay)

	assert.Equal(t, true, config.TradingOptions.TrailingStopOptions.Enable)
	assert.Equal(t, 0.4, config.TradingOptions.TrailingStopOptions.TrailingStopLoss)
	assert.Equal(t, 0.1, config.TradingOptions.TrailingStopOptions.TrailingTakeProfit)

	assert.Equal(t, true, len(config.TradingOptions.AllowList) > 0)
	assert.Contains(t, config.TradingOptions.AllowList, "AAVE")

	assert.Equal(t, true, len(config.TradingOptions.DenyList) > 0)
	assert.Contains(t, config.TradingOptions.DenyList, "GBPUSDT")
}
