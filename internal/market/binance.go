package market

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"strconv"
	"strings"
	"time"
)

// Ensures Binance implements the Market interface.
var _ Market = (*Binance)(nil)

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
		priceAsFloat, _ := strconv.ParseFloat(price.Price, 64)
		coin := Coin{
			Symbol: price.Symbol,
			Price:  priceAsFloat,
			Time:   now,
		}
		if coin.IsAvailableForTrading(b.config.TradingOptions.AllowList, b.config.TradingOptions.DenyList, b.config.TradingOptions.PairWith) {
			coins[coin.Symbol] = coin
		}
	}

	return coins, nil
}

func (b *Binance) GetSymbolInfo(ctx context.Context, symbol string) (SymbolInfo, error) {
	info, err := b.client.NewExchangeInfoService().Symbol(symbol).Do(ctx)
	if err != nil {
		return SymbolInfo{}, err
	}

	for _, s := range info.Symbols {
		if s.Symbol == strings.ToUpper(symbol) {
			stepSize := strings.Index(s.LotSizeFilter().StepSize, "1")
			if stepSize == -1 {
				stepSize = 0
			}

			return SymbolInfo{
				Symbol:   s.Symbol,
				StepSize: stepSize,
			}, nil
		}
	}

	return SymbolInfo{}, SymbolNotFoundError
}
