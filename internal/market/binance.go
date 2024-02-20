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

func (b *Binance) GetCoins(ctx context.Context) (CoinMap, error) {
	prices, err := b.client.NewListPricesService().Do(ctx)
	if err != nil {
		return nil, err
	}

	coins := make(CoinMap)
	now := time.Now()

	for _, price := range prices {
		coin := Coin{
			Symbol: price.Symbol,
			Price:  price.Price,
			Time:   now,
		}
		if coin.IsAvailableForTrading(b.config.TradingOptions.AllowList, b.config.TradingOptions.DenyList, b.config.TradingOptions.PairWith) {
			coins[coin.Symbol] = coin
		}
	}

	return coins, nil
}
