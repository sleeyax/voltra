package bot

import (
	"context"
	"fmt"
	"github.com/sleeyax/voltra/internal/config"
	"github.com/sleeyax/voltra/internal/database"
	"github.com/sleeyax/voltra/internal/database/models"
	"github.com/sleeyax/voltra/internal/market"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
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
	panic("implement me")
}

func (m *mockMarket) GetCoinsVolume(_ context.Context) (market.TradeVolumes, error) {
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
		StepSize: 0.0000001,
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
	m.orders[order.Symbol+string(order.Type)] = order
}

func (m *mockDatabase) HasOrder(orderType models.OrderType, market, symbol string) bool {
	order, ok := m.orders[symbol+string(orderType)]
	return ok && order.Type == orderType && order.Market == market
}

func (m *mockDatabase) CountOrders(orderType models.OrderType, market string) int64 {
	var count int64
	for _, order := range m.orders {
		if order.Type == orderType && order.Market == market {
			count++
		}
	}
	return count
}

func (m *mockDatabase) GetOrders(orderType models.OrderType, market string) []models.Order {
	var orders []models.Order
	for _, order := range m.orders {
		if order.Type == orderType && order.Market == market {
			orders = append(orders, order)
		}
	}
	return orders
}

func (m *mockDatabase) DeleteOrder(order models.Order) {
	delete(m.orders, order.Symbol+string(order.Type))
}

func (m *mockDatabase) SaveCache(_ models.Cache) {
	// ignore
}

func (m *mockDatabase) GetCache(_ string) (models.Cache, bool) {
	return models.Cache{}, false
}

func (m *mockDatabase) GetLastOrder(orderType models.OrderType, market, symbol string) (models.Order, bool) {
	panic("implement me")
}

func TestBot_buy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &config.Configuration{
		EnableTestMode: true,
		LoggingOptions: config.LoggingOptions{Enable: false},
		TradingOptions: config.TradingOptions{
			ChangeInPrice:        10, // 10%
			PairWith:             "USDT",
			Quantity:             10, // trade 10 USDT
			MinQuoteVolumeTraded: 100_000,
		},
	}

	m := newMockMarket(cancel)
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol:            "BTCUSDT",
			Price:             10_000,
			QuoteVolumeTraded: 50_000,
		},
		"ETH": market.Coin{
			Symbol:            "ETHUSDT",
			Price:             10_000,
			QuoteVolumeTraded: 80_000,
		},
	})
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol:            "BTCUSDT",
			Price:             10_500,
			QuoteVolumeTraded: 100_000,
		},
		"ETH": market.Coin{
			Symbol:            "ETHUSDT",
			Price:             9_000,
			QuoteVolumeTraded: 20_000,
		},
	})
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol:            "BTCUSDT",
			Price:             11_000,
			QuoteVolumeTraded: 120_000,
		},
		"ETH": market.Coin{
			Symbol:            "ETHUSDT",
			Price:             10_000,
			QuoteVolumeTraded: 400_000,
		},
	})
	// price change above threshold but not enough trading volume
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol:            "BTCUSDT",
			Price:             14_000,
			QuoteVolumeTraded: 30_000,
		},
		"ETH": market.Coin{
			Symbol:            "ETHUSDT",
			Price:             13_000,
			QuoteVolumeTraded: 40_000,
		},
	})

	db := newMockDatabase()

	b := New(c, m, db)

	var wg sync.WaitGroup
	wg.Add(1)
	b.buy(ctx, &wg)

	orders := db.GetOrders(models.BuyOrder, m.Name())
	assert.Equal(t, 1, len(orders))
	assert.Equal(t, "BTCUSDT", orders[0].Symbol)
	assert.Equal(t, 0.0009091, orders[0].Volume)
}

func TestBot_sell(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &config.Configuration{
		EnableTestMode: true,
		LoggingOptions: config.LoggingOptions{Enable: false},
		TradingOptions: config.TradingOptions{
			ChangeInPrice:   0.5,
			PairWith:        "USDT",
			Quantity:        15,
			TakeProfit:      0.1,
			StopLoss:        5,
			TradingFeeTaker: 0.075,
		},
	}

	m := newMockMarket(cancel)
	m.AddCoins(market.CoinMap{
		"XTZUSDT": market.Coin{
			Symbol: "XTZUSDT",
			Price:  1.295,
		},
	})

	db := newMockDatabase()
	db.SaveOrder(models.Order{
		Order: market.Order{
			Symbol: "XTZUSDT",
			Price:  1.292,
		},
		Market:     m.Name(),
		Type:       models.BuyOrder,
		Volume:     11.6,
		TakeProfit: &c.TradingOptions.TakeProfit,
		StopLoss:   &c.TradingOptions.StopLoss,
		IsTestMode: true,
	})

	b := New(c, m, db)

	var wg sync.WaitGroup
	wg.Add(1)
	b.sell(ctx, &wg)

	assert.Equal(t, int64(1), db.CountOrders(models.SellOrder, m.Name()))
	assert.Equal(t, int64(0), db.CountOrders(models.BuyOrder, m.Name()))
}

func TestBot_sell_with_trailing_stop_loss(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	c := &config.Configuration{
		EnableTestMode: true,
		LoggingOptions: config.LoggingOptions{Enable: false},
		TradingOptions: config.TradingOptions{
			ChangeInPrice:   0.5,
			PairWith:        "USDT",
			Quantity:        10,
			TakeProfit:      10,
			StopLoss:        5,
			TradingFeeTaker: 0.075, TrailingStopOptions: config.TrailingStopOptions{
				Enable:             true,
				TrailingStopLoss:   1,
				TrailingTakeProfit: 1,
			},
		},
	}

	m := newMockMarket(cancel)
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol: "BTC",
			Price:  11_000,
		},
	})
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol: "BTC",
			Price:  11_050,
		},
	})
	m.AddCoins(market.CoinMap{
		"BTC": market.Coin{
			Symbol: "BTC",
			Price:  9000,
		},
	})

	db := newMockDatabase()
	db.SaveOrder(models.Order{
		Order: market.Order{
			Symbol: "BTC",
			Price:  10_000,
		},
		Market:     m.Name(),
		Type:       models.BuyOrder,
		Volume:     0.000909,
		TakeProfit: &c.TradingOptions.TakeProfit,
		StopLoss:   &c.TradingOptions.StopLoss,
		IsTestMode: true,
	})

	b := New(c, m, db)

	var wg sync.WaitGroup
	wg.Add(1)
	b.sell(ctx, &wg)

	assert.Equal(t, int64(0), db.CountOrders(models.BuyOrder, m.Name()))
	orders := db.GetOrders(models.SellOrder, m.Name())
	assert.Equal(t, 1, len(orders))
	assert.NotNil(t, orders[0].PriceChangePercentage)
	assert.Equal(t, float64(-10), *orders[0].PriceChangePercentage)
	assert.Equal(t, float64(9000), orders[0].Price)
}

func TestBot_convertVolume(t *testing.T) {
	c := config.Configuration{
		EnableTestMode: true,
		LoggingOptions: config.LoggingOptions{Enable: false},
	}
	b := New(&c, newMockMarket(nil), newMockDatabase())
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
	assert.Equal(t, 0.0009091, v)
}
