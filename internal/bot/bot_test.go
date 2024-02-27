package bot

import (
	"context"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockMarket struct{}

// ensure mockMarket implements the Market interface
var _ market.Market = (*mockMarket)(nil)

func (m mockMarket) Name() string {
	return "mock market"
}

func (m mockMarket) Buy(ctx context.Context, coin string, quantity float64) (market.Order, error) {
	panic("implement me")
}

func (m mockMarket) Sell(ctx context.Context, coin string, quantity float64) (market.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockMarket) GetCoins(_ context.Context) (market.CoinMap, error) {
	panic("implement me")
}

func (m mockMarket) GetSymbolInfo(_ context.Context, symbol string) (market.SymbolInfo, error) {
	return market.SymbolInfo{
		Symbol:   "BTC",
		StepSize: 6,
	}, nil
}

func TestBot_convertVolume(t *testing.T) {
	c := config.Configuration{}
	b := New(&c, mockMarket{}, nil)
	v, err := b.convertVolume(context.Background(), 50, market.VolatileCoin{
		Coin: market.Coin{
			Symbol: "BTC",
			Price:  100,
		},
		Percentage: 15,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0.5, v)

	v, err = b.convertVolume(context.Background(), 50, market.VolatileCoin{
		Coin: market.Coin{
			Symbol: "BTC",
			Price:  10000,
		},
		Percentage: 15,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0.005, v)
}
