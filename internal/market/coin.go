package market

import (
	"fmt"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/utils"
	"time"
)

type Coin struct {
	// The symbol of the coin.
	Symbol string `json:"symbol"`

	// The price of the coin.
	Price string `json:"price"`

	// The time this coin was indexed.
	Time time.Time `json:"time"`
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
