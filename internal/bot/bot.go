package bot

import (
	"context"
	"fmt"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database/models"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/utils"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Bot struct {
	market  market.Market
	db      database.Database
	history *History
	config  *config.Configuration
	log     *zap.SugaredLogger
}

func New(config *config.Configuration, market market.Market, db database.Database) *Bot {
	var logger *zap.Logger
	if config.ScriptOptions.EnableStructuredLogging {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	sugar := logger.Sugar()
	history := NewHistory(config.TradingOptions.RecheckInterval)
	return &Bot{market: market, history: history, config: config, log: sugar, db: db}
}

func (b *Bot) Close() error {
	return b.log.Sync()
}

// Monitor starts monitoring the market for price changes.
func (b *Bot) Monitor(ctx context.Context) error {
	// Seed initial data on bot startup.
	b.log.Info("Bot started.")
	if err := b.updateLatestCoins(ctx); err != nil {
		return err
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

		if b.db.HasOrder(models.BuyOrder, b.market.Name(), volatileCoin.Symbol) {
			b.log.Warnf("Already bought %s. Skipping.", volatileCoin.Symbol)
			continue
		}

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

		// Pretend to buy the coin and save the order if test mode is enabled.
		if b.config.ScriptOptions.TestMode {
			b.db.SaveOrder(models.Order{
				BuyOrder: market.BuyOrder{
					OrderID:         0,
					Symbol:          volatileCoin.Symbol,
					Price:           volatileCoin.Price,
					TransactionTime: time.Now(),
				},
				Market:     b.market.Name(),
				Type:       models.BuyOrder,
				Volume:     volume,
				IsTestMode: true,
			})
			continue
		}

		// Otherwise, buy the coin and save the real order.
		buyOrder, err := b.market.Buy(ctx, volatileCoin.Symbol, volume)
		if err != nil {
			return err
		}
		b.db.SaveOrder(models.Order{
			BuyOrder: buyOrder,
			Market:   b.market.Name(),
			Type:     models.BuyOrder,
			Volume:   volume,
		})
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
