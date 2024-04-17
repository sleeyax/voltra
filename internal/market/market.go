package market

import (
	"context"
	"errors"
)

var SymbolNotFoundError = errors.New("symbol not found")

type Market interface {
	// Name returns the name of the market.
	Name() string

	// GetCoins returns the current price of all coins on the market.
	GetCoins(ctx context.Context) (Coins, error)

	// GetCoinsVolume returns the quote volume traded for all coins on the market.
	GetCoinsVolume(ctx context.Context) (TradeVolumes, error)

	// GetSymbolInfo returns the symbol info for the given symbol.
	GetSymbolInfo(ctx context.Context, symbol string) (SymbolInfo, error)

	// Buy buys the given quantity of the given coin.
	Buy(ctx context.Context, coin string, quantity float64) (Order, error)

	// Sell sells the given quantity of the given coin.
	Sell(ctx context.Context, coin string, quantity float64) (Order, error)
}
