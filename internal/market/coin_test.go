package market

import (
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestCoin_IsAvailableForTrading(t *testing.T) {
	allowList := []string{"BTC"}
	denyList := []string{"EURUSDT"}
	pairWith := "USDT"
	coin := Coin{
		Price: "10000",
		Time:  time.Now(),
	}

	// test allowlist and denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, coin.IsAvailableForTrading(allowList, denyList, pairWith), true)
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, coin.IsAvailableForTrading(allowList, denyList, pairWith), false)
	coin.Symbol = "EURUSDT"
	assert.Equal(t, coin.IsAvailableForTrading(allowList, denyList, pairWith), false)

	// test no allowlist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, coin.IsAvailableForTrading([]string{}, denyList, pairWith), true)
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, coin.IsAvailableForTrading([]string{}, denyList, pairWith), true)
	coin.Symbol = "EURUSDT"
	assert.Equal(t, coin.IsAvailableForTrading([]string{}, denyList, pairWith), false)

	// test no denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, coin.IsAvailableForTrading(allowList, []string{}, pairWith), true)
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, coin.IsAvailableForTrading(allowList, []string{}, pairWith), false)
	coin.Symbol = "EURUSDT"
	assert.Equal(t, coin.IsAvailableForTrading(allowList, []string{}, pairWith), false)

	// test no allowlist and no denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, coin.IsAvailableForTrading([]string{}, []string{}, pairWith), true)
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, coin.IsAvailableForTrading([]string{}, []string{}, pairWith), true)
	coin.Symbol = "EURUSDT"
	assert.Equal(t, coin.IsAvailableForTrading([]string{}, []string{}, pairWith), true)
}
