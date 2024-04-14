package market

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCoin_IsAvailableForTrading(t *testing.T) {
	allowList := []string{"BTC"}
	denyList := []string{"EURUSDT"}
	pairWith := "USDT"
	coin := Coin{
		Price: 10000,
		Time:  time.Now(),
	}

	// test allowlist and denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading(allowList, denyList, pairWith))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith))

	// test no allowlist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, denyList, pairWith))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, denyList, pairWith))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading([]string{}, denyList, pairWith))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading([]string{}, denyList, pairWith))

	// test no denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading(allowList, []string{}, pairWith))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, []string{}, pairWith))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, []string{}, pairWith))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, []string{}, pairWith))

	// test no allowlist and no denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, []string{}, pairWith))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, []string{}, pairWith))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, []string{}, pairWith))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading([]string{}, []string{}, pairWith))
}
