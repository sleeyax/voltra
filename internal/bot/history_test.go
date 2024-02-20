package bot

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHistory_Size(t *testing.T) {
	history := NewHistory(3)
	assert.Equal(t, 0, history.Size())
	history.AddRecord(nil)
	assert.Equal(t, 1, history.Size())
}

func TestHistory_Min(t *testing.T) {
	history := NewHistory(3)
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 10000.0},
	})
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 8000.0},
	})
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 20000.0},
	})
	m := history.Min("BTCUSDT")
	assert.Equal(t, 8000.0, m.coins["BTCUSDT"].Price)
}

func TestHistory_Max(t *testing.T) {
	history := NewHistory(3)
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 10000.0},
	})
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 20000.0},
	})
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 8000.0},
	})
	m := history.Max("BTCUSDT")
	assert.Equal(t, 20000.0, m.coins["BTCUSDT"].Price)
}

func TestHistory_IdentifyVolatileCoins(t *testing.T) {
	// Basic percentage increase check.
	history := NewHistory(UnlimitedHistoryLength)
	percentage := 15.0
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 20_000.0},
	})
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 23_000.0},
	})
	v := history.IdentifyVolatileCoins(percentage)
	assert.Equal(t, percentage, v["BTCUSDT"])

	// This coin is already identified as volatile above.
	// a sudden spike in price shouldn't affect the result within the current time window.
	history.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 25_000.0},
	})
	v = history.IdentifyVolatileCoins(percentage)
	assert.Equal(t, percentage, v["BTCUSDT"])

	// A brand-new coin should not yet be volatile.
	history.AddRecord(market.CoinMap{
		"ETHUSDT": {Price: 3000},
	})
	v = history.IdentifyVolatileCoins(percentage)
	vv, ok := v["ETHUSDT"]
	assert.Equal(t, false, ok, vv)

	// Test price drop.
	history.AddRecord(market.CoinMap{
		"ETHUSDT": {Price: 2000},
	})
	v = history.IdentifyVolatileCoins(percentage)
	vv, ok = v["ETHUSDT"]
	assert.Equal(t, false, ok, vv)

	// Test price increase
	history.AddRecord(market.CoinMap{
		"ETHUSDT": {Price: 10_000},
	})
	v = history.IdentifyVolatileCoins(percentage)
	vv, ok = v["ETHUSDT"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 400.0, vv)
}
