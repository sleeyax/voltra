package market

import (
	"context"
	"github.com/adshao/go-binance/v2"
	"github.com/sleeyax/voltra/internal/config"
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
	client := binance.NewClient(m.AccessKey, m.SecretKey)
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
		coins[coin.Symbol] = coin
	}

	return coins, nil
}

func (b *Binance) GetCoinsVolume(ctx context.Context) (TradeVolumes, error) {
	volumeMap := make(TradeVolumes)
	if b.config.TradingOptions.MinQuoteVolumeTraded != 0.0 {
		priceStats24Hours, err := b.client.NewListPriceChangeStatsService().Do(ctx)
		if err != nil {
			return nil, err
		}
		for _, priceStat := range priceStats24Hours {
			quoteVolumeAsFloat, _ := strconv.ParseFloat(priceStat.QuoteVolume, 64)
			volumeMap[priceStat.Symbol] = quoteVolumeAsFloat
		}
	}
	return volumeMap, nil
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

	marketOrder, err := b.client.NewCreateOrderService().
		Symbol(coin).
		Side(side).
		Type(binance.OrderTypeMarket).
		Quantity(quantityAsString).
		Do(ctx)

	if err != nil {
		return Order{}, err
	}

	order := Order{
		OrderID:         marketOrder.OrderID,
		Symbol:          marketOrder.Symbol,
		TransactionTime: time.Unix(marketOrder.TransactTime, 0),
	}

	// Market orders are not always filled at one singular price.
	// If that's the case, we need to find the averages of all 'parts' (fills) of this order in order to calculate the total price (see code below).
	// Otherwise, it's safe to read the price from the order itself.
	if len(marketOrder.Fills) == 0 {
		p, _ := strconv.ParseFloat(marketOrder.Price, 64)
		order.Price = p
		return order, nil
	}

	// Calculate the average price of all fills.
	var totalPrice float64
	var totalQuantity float64

	for _, fill := range marketOrder.Fills {
		qty, _ := strconv.ParseFloat(fill.Quantity, 64)
		price, _ := strconv.ParseFloat(fill.Price, 64)
		totalQuantity += qty
		totalPrice += price * qty
	}

	fillAvg := totalPrice / totalQuantity

	order.Price = fillAvg

	return order, nil
}

func (b *Binance) Buy(ctx context.Context, coin string, quantity float64) (Order, error) {
	return b.executeOrder(ctx, coin, quantity, binance.SideTypeBuy)
}

func (b *Binance) Sell(ctx context.Context, coin string, quantity float64) (Order, error) {
	return b.executeOrder(ctx, coin, quantity, binance.SideTypeSell)
}
