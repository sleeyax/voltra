package bot

import (
	"github.com/sleeyax/voltra/internal/market"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVolatilityWindow_Size(t *testing.T) {
	window := NewVolatilityWindow(3)
	assert.Equal(t, 0, window.Size())
	window.AddRecord(nil)
	assert.Equal(t, 1, window.Size())
}

func TestVolatilityWindow_Min(t *testing.T) {
	window := NewVolatilityWindow(3)
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 10000.0},
	})
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 8000.0},
	})
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 20000.0},
	})
	m := window.Min("BTCUSDT")
	assert.Equal(t, 8000.0, m.coins["BTCUSDT"].Price)
}

func TestVolatilityWindow_Max(t *testing.T) {
	window := NewVolatilityWindow(3)
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 10000.0},
	})
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 20000.0},
	})
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 8000.0},
	})
	m := window.Max("BTCUSDT")
	assert.Equal(t, 20000.0, m.coins["BTCUSDT"].Price)
}

func TestVolatilityWindow_IdentifyVolatileCoins(t *testing.T) {
	// Basic percentage increase check.
	window := NewVolatilityWindow(UnlimitedVolatilityWindowLength)
	percentage := 15.0
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 20_000.0},
	})
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 23_000.0},
	})
	v := window.IdentifyVolatileCoins(percentage)
	assert.Equal(t, percentage, v["BTCUSDT"].Percentage)

	// This coin is already identified as volatile above.
	// a sudden spike in price shouldn't affect the result within the current time window.
	window.AddRecord(market.CoinMap{
		"BTCUSDT": {Price: 25_000.0},
	})
	v = window.IdentifyVolatileCoins(percentage)
	assert.Equal(t, percentage, v["BTCUSDT"].Percentage)

	// A brand-new coin should not yet be volatile.
	window.AddRecord(market.CoinMap{
		"ETHUSDT": {Price: 3000},
	})
	v = window.IdentifyVolatileCoins(percentage)
	vv, ok := v["ETHUSDT"]
	assert.Equal(t, false, ok, vv.Percentage)

	// Test price drop.
	window.AddRecord(market.CoinMap{
		"ETHUSDT": {Price: 2000},
	})
	v = window.IdentifyVolatileCoins(percentage)
	vv, ok = v["ETHUSDT"]
	assert.Equal(t, false, ok, vv.Percentage)

	// Test price increase
	window.AddRecord(market.CoinMap{
		"ETHUSDT": {Price: 10_000},
	})
	v = window.IdentifyVolatileCoins(percentage)
	vv, ok = v["ETHUSDT"]
	assert.Equal(t, true, ok)
	assert.Equal(t, 400.0, vv.Percentage)
}
