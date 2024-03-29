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
// It checks if the coin is in the custom list and if it's not a fiat currency, both of which are options defined in the given config.
func (c Coin) IsAvailableForTrading(allowList, denyList []string, pairWith string) bool {
	if len(allowList) > 0 {
		if !utils.Any(allowList, func(allowedCoin string) bool {
			allowedSymbol := fmt.Sprintf("%s%s", allowedCoin, pairWith)
			return allowedSymbol == c.Symbol
		}) {
			return false
		}
	}

	if len(denyList) > 0 {
		if utils.Any(denyList, func(fiat string) bool {
			return c.Symbol == fiat
		}) {
			return false
		}
	}

	return true
}
