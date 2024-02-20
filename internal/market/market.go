package market

import (
	"context"
)

type Market interface {
	// GetCoins returns the current price of all coins on the market.
	GetCoins(ctx context.Context) (CoinMap, error)
}
