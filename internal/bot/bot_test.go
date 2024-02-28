package bot

import (
	"context"
	"fmt"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database/models"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"testing"
	"time"
)

type mockMarket struct {
	coinsIndex int
	coins      []market.CoinMap
	cancel     context.CancelFunc
}

// ensure mockMarket implements the Market interface
var _ market.Market = (*mockMarket)(nil)

func newMockMarket(cancel context.CancelFunc) *mockMarket {
	return &mockMarket{
		coins:  make([]market.CoinMap, 0),
		cancel: cancel,
	}
}

func (m *mockMarket) Name() string {
	return "mock market"
}

func (m *mockMarket) Buy(ctx context.Context, coin string, quantity float64) (market.Order, error) {
	panic("implement me")
}

func (m *mockMarket) Sell(ctx context.Context, coin string, quantity float64) (market.Order, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockMarket) GetCoins(_ context.Context) (market.CoinMap, error) {
	if m.coinsIndex >= len(m.coins) {
		m.cancel()
		return nil, fmt.Errorf("no coins found at index %d", m.coinsIndex)
	}

	coins := m.coins[m.coinsIndex]

	m.coinsIndex++

	return coins, nil
}

func (m *mockMarket) AddCoins(coins market.CoinMap) {
	m.coins = append(m.coins, coins)
}

func (m *mockMarket) GetSymbolInfo(_ context.Context, symbol string) (market.SymbolInfo, error) {
	return market.SymbolInfo{
		Symbol:   "BTC",
		StepSize: 6,
	}, nil
}

type mockDatabase struct {
	orders map[string]models.Order
}

var _ database.Database = (*mockDatabase)(nil)

func newMockDatabase() *mockDatabase {
	return &mockDatabase{
		orders: make(map[string]models.Order),
	}
}

func (m *mockDatabase) SaveOrder(order models.Order) {
	m.orders[order.Symbol] = order
}

func (m *mockDatabase) HasOrder(orderType models.OrderType, market, symbol string) bool {
	_, ok := m.orders[symbol]
	return ok
}

func (m *mockDatabase) CountOrders(orderType models.OrderType, market string) int64 {
	return int64(len(m.orders))
}

func (m *mockDatabase) GetOrders(orderType models.OrderType, market string) []models.Order {
	return maps.Values(m.orders)
}

func (m *mockDatabase) DeleteOrder(order models.Order) {
	delete(m.orders, order.Symbol)
}

func TestBot_buy_volatile_coin(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)

	c := &config.Configuration{
		ScriptOptions: config.ScriptOptions{
			TestMode:       true,
			DisableLogging: true,
		},
		TradingOptions: config.TradingOptions{
			ChangeInPrice: 10, // 10%
			PairWith:      "USDT",
			Quantity:      10, // trade 10 USDT
		},
	}

	m := newMockMarket(cancel)
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol: "BTC",
			Price:  10_000,
		},
	})
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol: "BTC",
			Price:  11_000,
		},
	})

	db := newMockDatabase()

	b := New(c, m, db)

	b.buy(ctx)

	orders := db.GetOrders(models.BuyOrder, m.Name())
	assert.Equal(t, 1, len(orders))
	assert.Equal(t, "BTC", orders[0].Symbol)
	assert.Equal(t, 0.000909, orders[0].Volume)
}

func TestBot_convertVolume(t *testing.T) {
	c := config.Configuration{
		ScriptOptions: config.ScriptOptions{
			TestMode:       true,
			DisableLogging: true,
		},
	}
	b := New(&c, newMockMarket(nil), nil)
	v, err := b.convertVolume(context.Background(), 50, market.VolatileCoin{
		Coin: market.Coin{
			Symbol: "BTC",
			Price:  100,
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0.5, v)

	v, err = b.convertVolume(context.Background(), 50, market.VolatileCoin{
		Coin: market.Coin{
			Symbol: "BTC",
			Price:  10000,
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0.005, v)

	v, err = b.convertVolume(context.Background(), 10, market.VolatileCoin{
		Coin: market.Coin{
			Symbol: "BTC",
			Price:  11_000,
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0.000909, v)
}
