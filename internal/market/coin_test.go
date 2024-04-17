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
	minQuoteVolumeTraded := 5000.0
	coin := Coin{
		Price:             10000,
		QuoteVolumeTraded: 40000,
		Time:              time.Now(),
	}

	// test allowlist and denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))

	// test no allowlist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading([]string{}, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading([]string{}, denyList, pairWith, minQuoteVolumeTraded))

	// test no denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading(allowList, []string{}, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, []string{}, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, []string{}, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, []string{}, pairWith, minQuoteVolumeTraded))

	// test no allowlist and no denylist
	coin.Symbol = "BTCUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, []string{}, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "ETHUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, []string{}, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "EURUSDT"
	assert.Equal(t, true, coin.IsAvailableForTrading([]string{}, []string{}, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "BTCUSDC"
	assert.Equal(t, false, coin.IsAvailableForTrading([]string{}, []string{}, pairWith, minQuoteVolumeTraded))

	// test 24H quote volume threshold
	coin.Symbol = "BTCUSDT"
	coin.QuoteVolumeTraded = 4000
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "BTCUSDT"
	coin.QuoteVolumeTraded = 5000
	assert.Equal(t, true, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "BTCUSDT"
	coin.QuoteVolumeTraded = 5000
	assert.Equal(t, true, coin.IsAvailableForTrading(allowList, denyList, pairWith, 5000))
	coin.Symbol = "BTCUSDT"
	coin.QuoteVolumeTraded = 6000.0
	assert.Equal(t, true, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "ETHUSDT"
	coin.QuoteVolumeTraded = 6000
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
	coin.Symbol = "BTCUSDC"
	coin.QuoteVolumeTraded = 10000
	assert.Equal(t, false, coin.IsAvailableForTrading(allowList, denyList, pairWith, minQuoteVolumeTraded))
}
