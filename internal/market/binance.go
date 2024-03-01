package market

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"strconv"
	"time"
)

// Ensures Binance implements the Market interface.
var _ Market = (*Binance)(nil)

type Binance struct {
	config config.Configuration
	client *binance.Client
}

func NewBinance(config config.Configuration) *Binance {
	m := config.Markets.Binance
	client := binance.NewClient(m.ApiKey, m.SecretKey)
	return &Binance{config: config, client: client}
}

func (b *Binance) Name() string {
	return "binance"
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
		if s.Symbol == symbol {
			stepSize, _ := strconv.ParseFloat(s.LotSizeFilter().StepSize, 64)

			return SymbolInfo{
				Symbol:   s.Symbol,
				StepSize: stepSize,
			}, nil
		}
	}

	return SymbolInfo{}, SymbolNotFoundError
}

func (b *Binance) executeOrder(ctx context.Context, coin string, quantity float64, side binance.SideType) (Order, error) {
	quantityAsString := strconv.FormatFloat(quantity, 'f', -1, 64)

	order, err := b.client.NewCreateOrderService().
		Symbol(coin).
		Side(side).
		Type(binance.OrderTypeMarket).
		Quantity(quantityAsString).
		Do(ctx)

	if err != nil {
		return Order{}, err
	}

	p, _ := strconv.ParseFloat(order.Price, 64)

	return Order{
		OrderID:         order.OrderID,
		Symbol:          order.Symbol,
		Price:           p,
		TransactionTime: time.Unix(order.TransactTime, 0),
	}, err
}

func (b *Binance) Buy(ctx context.Context, coin string, quantity float64) (Order, error) {
	return b.executeOrder(ctx, coin, quantity, binance.SideTypeBuy)
}

func (b *Binance) Sell(ctx context.Context, coin string, quantity float64) (Order, error) {
	return b.executeOrder(ctx, coin, quantity, binance.SideTypeSell)
}
