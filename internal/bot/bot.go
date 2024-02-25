package bot

import (
	"context"
	"fmt"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/utils"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Bot struct {
	market    market.Market
	history   *History
	config    config.Configuration
	log       *zap.SugaredLogger
	buyOrders market.BuyOrderMap
}

func New(c config.Configuration, m market.Market) *Bot {
	var logger *zap.Logger
	if c.ScriptOptions.EnableStructuredLogging {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	sugar := logger.Sugar()
	return &Bot{market: m, history: NewHistory(c.TradingOptions.RecheckInterval), config: c, log: sugar, buyOrders: make(market.BuyOrderMap)}
}

func (b *Bot) Close() error {
	return b.log.Sync()
}

// Monitor starts monitoring the market for price changes.
func (b *Bot) Monitor(ctx context.Context) error {
	// Seed initial data on bot startup.
	if b.history.Size() == 0 {
		b.log.Info("Bot started.")
		if err := b.updateLatestCoins(ctx); err != nil {
			return err
		}
	}

	// Sleep until the next recheck interval.
	lastRecord := b.history.GetLatestRecord()
	delta := utils.CalculateTimeDelta(b.config.TradingOptions.TimeDifference, b.config.TradingOptions.RecheckInterval)
	if time.Since(lastRecord.time) < delta {
		interval := delta - time.Since(lastRecord.time)
		b.log.Infof("Sleeping %s.", interval)
		time.Sleep(interval)
	}

	if err := b.updateLatestCoins(ctx); err != nil {
		return err
	}

	// TODO: Implement trading logic.
	volatileCoins := b.history.IdentifyVolatileCoins(b.config.TradingOptions.ChangeInPrice)
	for _, volatileCoin := range volatileCoins {
		b.log.Infof("Coin %s has gained %f%% within the last %d minutes.", volatileCoin.Symbol, volatileCoin.Percentage, b.config.TradingOptions.TimeDifference)

		volume, err := b.ConvertVolume(ctx, b.config.TradingOptions.Quantity, volatileCoin)
		if err != nil {
			return err
		}

		b.log.Infow(fmt.Sprintf("Trading %f %s of %s.", volume, b.config.TradingOptions.PairWith, volatileCoin.Symbol),
			"volume", volume,
			"pair_with", b.config.TradingOptions.PairWith,
			"symbol", volatileCoin.Symbol,
			"price", volatileCoin.Price,
			"percentage", volatileCoin.Percentage,
			"testMode", b.config.ScriptOptions.TestMode,
		)
		if err = b.Buy(ctx, volume, volatileCoin, b.config.ScriptOptions.TestMode); err != nil {
			return err
		}
		// TODO: save buy orders to local database
	}

	return nil
}

// updateLatestCoins fetches the latest coins from the market and appends them to the history.
func (b *Bot) updateLatestCoins(ctx context.Context) error {
	b.log.Info("Fetching latest coins.")

	coins, err := b.market.GetCoins(ctx)
	if err != nil {
		return err
	}

	b.history.AddRecord(coins)

	return nil
}

// ConvertVolume converts the volume given in the configured quantity from base currency (USDT) to each coin's volume.
func (b *Bot) ConvertVolume(ctx context.Context, quantity float64, volatileCoin market.VolatileCoin) (float64, error) {
	info, err := b.market.GetSymbolInfo(ctx, volatileCoin.Symbol)
	if err != nil {
		return 0, err
	}

	volume := quantity / volatileCoin.Price

	// Round the volume to the step size of the coin.
	if info.StepSize != 0 {
		formattedVolumeString := strconv.FormatFloat(volume, 'f', info.StepSize, 64)
		volume, err = strconv.ParseFloat(formattedVolumeString, 64)
		if err != nil {
			return 0, err
		}
	}

	return volume, nil
}

func (b *Bot) Buy(ctx context.Context, volume float64, volatileCoin market.VolatileCoin, isTestMode bool) error {
	_, ok := b.buyOrders[volatileCoin.Symbol]
	if ok {
		b.log.Infof("Already bought %s. Skipping.", volatileCoin.Symbol)
		return nil
	}

	if isTestMode {
		b.buyOrders[volatileCoin.Symbol] = market.BuyOrder{
			OrderID:          0,
			Symbol:           volatileCoin.Symbol,
			Price:            volatileCoin.Price,
			TransactionTime:  time.Now().Unix(),
			ExecutedQuantity: strconv.FormatFloat(volume, 'f', -1, 64),
		}
		return nil
	}

	buyOrder, err := b.market.Buy(ctx, volatileCoin.Symbol, volume)
	if err != nil {
		return err
	}

	b.buyOrders[volatileCoin.Symbol] = buyOrder

	return nil
}
