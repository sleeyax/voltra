package market

import (
	"fmt"
	"github.com/sleeyax/voltra/internal/utils"
	"time"
)

type Coin struct {
	// The symbol of the coin.
	Symbol string `json:"symbol"`

	// The price of the coin.
	Price float64 `json:"price"`

	// The 24h quote asset volume traded of the coin.
	QuoteVolumeTraded float64 `json:"quote_volume_traded"`

	// The time this coin was indexed.
	Time time.Time `json:"time"`
}

type VolatileCoin struct {
	// The coin that has gained in price.
	Coin

	// Percentage of price increase.
	Percentage float64
}

type VolatileCoins map[string]VolatileCoin

type CoinMap map[string]Coin

type CoinVolumeTradedMap map[string]float64

var CoinVolumes CoinVolumeTradedMap

type SymbolInfo struct {
	// The symbol of the coin.
	Symbol string

	// The step size of the coin.
	// E.g. 0.001.
	StepSize float64
}

func (c Coin) String() string {
	return c.Symbol
}

// IsAvailableForTrading checks if the coin should be picked up by the bot for trading.
// It checks whether the coin has the desired minimum quote asset trading volume, is in the custom list, and it's not a blacklisted symbol. These options are defined in the given config file.
func (c Coin) IsAvailableForTrading(allowList, denyList []string, pairWith string, minQuoteVolumeTraded float64) bool {
	if c.QuoteVolumeTraded < minQuoteVolumeTraded {
		return false
	}
	if len(denyList) > 0 {
		if utils.Any(denyList, func(blacklistedSymbol string) bool {
			return c.Symbol == blacklistedSymbol
		}) {
			return false
		}
	}

	if len(allowList) > 0 {
		if !utils.Any(allowList, func(allowedCoin string) bool {
			allowedSymbol := fmt.Sprintf("%s%s", allowedCoin, pairWith)
			return allowedSymbol == c.Symbol
		}) {
			return false
		}
	} else {
		return c.Symbol[len(c.Symbol)-len(pairWith):] == pairWith
	}

	return true
}
