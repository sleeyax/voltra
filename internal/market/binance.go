package market

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"time"
)

type Binance struct {
	config config.Configuration
	client *binance.Client
}

func NewBinance(config config.Configuration) *Binance {
	m := config.ScriptOptions.Markets.Binance
	client := binance.NewClient(m.ApiKey, m.SecretKey)
	return &Binance{config: config, client: client}
}

func (b *Binance) GetCoins(ctx context.Context) ([]Coin, error) {
	prices, err := b.client.NewListPricesService().Do(ctx)
	if err != nil {
		return nil, err
	}

	var coins []Coin

	for _, price := range prices {
		coin := Coin{
			Symbol: price.Symbol,
			Price:  price.Price,
			Time:   time.Now(),
		}
		if coin.IsAvailableForTrading(b.config.TradingOptions.AllowList, b.config.TradingOptions.DenyList, b.config.TradingOptions.PairWith) {
			coins = append(coins, coin)
		}
	}

	return coins, nil
}
