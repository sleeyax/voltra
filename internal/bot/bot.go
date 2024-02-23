package bot

import (
	"context"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/utils"
	"go.uber.org/zap"
	"time"
)

type Bot struct {
	market  market.Market
	history *History
	config  config.Configuration
	log     *zap.SugaredLogger
}

func New(c config.Configuration, m market.Market) *Bot {
	var logger *zap.Logger
	if c.ScriptOptions.EnableStructuredLogging {
		logger, _ = zap.NewProduction()
	} else {
		logger, _ = zap.NewDevelopment()
	}
	sugar := logger.Sugar()
	return &Bot{market: m, history: NewHistory(c.TradingOptions.RecheckInterval), config: c, log: sugar}
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
	for coin, percentage := range volatileCoins {
		b.log.Infof("Coin %s has gained %f%% within the last %d minutes.", coin, percentage, b.config.TradingOptions.TimeDifference)
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
