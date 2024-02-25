package bot

import (
	"context"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockMarket struct{}

func (m mockMarket) GetCoins(_ context.Context) (market.CoinMap, error) {
	//TODO implement me
	panic("implement me")
}

func (m mockMarket) GetSymbolInfo(_ context.Context, symbol string) (market.SymbolInfo, error) {
	return market.SymbolInfo{
		Symbol:   "BTC",
		StepSize: 6,
	}, nil
}

func TestBot_ConvertVolume(t *testing.T) {
	c := config.Configuration{}
	b := New(c, mockMarket{})
	v, err := b.ConvertVolume(context.Background(), 50, market.VolatileCoin{
		Coin: market.Coin{
			Symbol: "BTC",
			Price:  100,
		},
		Percentage: 15,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0.5, v)

	v, err = b.ConvertVolume(context.Background(), 50, market.VolatileCoin{
		Coin: market.Coin{
			Symbol: "BTC",
			Price:  10000,
		},
		Percentage: 15,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0.005, v)
}
