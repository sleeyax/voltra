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
	"go.uber.org/zap/zapcore"
	"math"
	"strconv"
	"time"
)

type Bot struct {
	market           market.Market
	db               database.Database
	volatilityWindow *VolatilityWindow
	config           *config.Configuration
	log              *zap.SugaredLogger
}

func New(config *config.Configuration, market market.Market, db database.Database) *Bot {
	var logger *zap.Logger
	if config.ScriptOptions.DisableLogging {
		logger = zap.NewNop()
	} else if config.ScriptOptions.EnableStructuredLogging {
		logger, _ = zap.NewProduction()
	} else {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		logger, _ = loggerConfig.Build()
	}
	sugar := logger.Sugar()
	window := NewVolatilityWindow(config.TradingOptions.RecheckInterval)
	return &Bot{market: market, volatilityWindow: window, config: config, log: sugar, db: db}
}

func (b *Bot) Close() error {
	return b.log.Sync()
}

// Start starts monitoring the market for price changes.
func (b *Bot) Start(ctx context.Context) {
	b.log.Info("Bot started.")
	go b.sell(ctx)
	b.buy(ctx)
}

func (b *Bot) buy(ctx context.Context) {
	b.log.Info("Watching coins to buy.")

	if err := b.updateLatestCoins(ctx); err != nil {
		panic(fmt.Sprintf("failed to load initial latest coins: %s", err))
	}

	for {
		select {
		case <-ctx.Done():
			b.log.Info("Bot stopped buying coins.")
			return
		default:
			// Wait until the next recheck interval.
			lastRecord := b.volatilityWindow.GetLatestRecord()
			delta := utils.CalculateTimeDuration(b.config.TradingOptions.TimeDifference, b.config.TradingOptions.RecheckInterval)
			if time.Since(lastRecord.time) < delta {
				interval := delta - time.Since(lastRecord.time)
				b.log.Infof("Sleeping %s.", interval)
				time.Sleep(interval)
			}

			// Fetch the latest coins again after the waiting period.
			if err := b.updateLatestCoins(ctx); err != nil {
				b.log.Errorf("Failed to update latest coins: %s.", err)
				continue
			}

			// Skip if the max amount of buy orders has been reached.
			if maxBuyOrders := int64(b.config.TradingOptions.MaxCoins); maxBuyOrders != 0 && b.db.CountOrders(models.BuyOrder, b.market.Name()) >= maxBuyOrders {
				b.log.Warnf("Max amount of buy orders reached.")
				continue
			}

			// Identify volatile coins in the current time window and trade them if any are found.
			volatileCoins := b.volatilityWindow.IdentifyVolatileCoins(b.config.TradingOptions.ChangeInPrice)
			b.log.Infof("Found %d volatile coins.", len(volatileCoins))
			for _, volatileCoin := range volatileCoins {
				b.log.Infof("Coin %s has gained %f%% within the last %d minutes.", volatileCoin.Symbol, volatileCoin.Percentage, b.config.TradingOptions.TimeDifference)

				// Skip if the coin has already been bought.
				if b.db.HasOrder(models.BuyOrder, b.market.Name(), volatileCoin.Symbol) {
					b.log.Warnf("Already bought %s. Skipping.", volatileCoin.Symbol)
					continue
				}

				// Determine the correct volume to buy based on the configured quantity.
				volume, err := b.convertVolume(ctx, b.config.TradingOptions.Quantity, volatileCoin)
				if err != nil {
					b.log.Errorf("Failed to convert volume: %s. Skipping the trade.", err)
					continue
				}

				b.log.Infow(fmt.Sprintf("Trading %f %s of %s.", volume, b.config.TradingOptions.PairWith, volatileCoin.Symbol),
					"volume", volume,
					"pair_with", b.config.TradingOptions.PairWith,
					"symbol", volatileCoin.Symbol,
					"price", volatileCoin.Price,
					"percentage", volatileCoin.Percentage,
					"testMode", b.config.ScriptOptions.TestMode,
				)

				order := models.Order{
					Market:     b.market.Name(),
					Type:       models.BuyOrder,
					Volume:     volume,
					TakeProfit: b.config.TradingOptions.TakeProfit,
					StopLoss:   b.config.TradingOptions.StopLoss,
				}

				// Pretend to buy the coin and save the order if test mode is enabled.
				if b.config.ScriptOptions.TestMode {
					order.Order = market.Order{
						OrderID:         0,
						Symbol:          volatileCoin.Symbol,
						Price:           volatileCoin.Price,
						TransactionTime: time.Now(),
					}
					order.IsTestMode = true
				} else {
					// Otherwise, buy the coin and save the real order.
					buyOrder, err := b.market.Buy(ctx, volatileCoin.Symbol, volume)
					if err != nil {
						b.log.Errorf("Failed to buy %s: %s.", volatileCoin.Symbol, err)
						continue
					}

					order.Order = buyOrder
				}

				b.db.SaveOrder(order)
			}
		}
	}
}

func (b *Bot) sell(ctx context.Context) {
	b.log.Info("Watching coins to sell.")

	for {
		select {
		case <-ctx.Done():
			b.log.Info("Bot stopped selling coins.")
			return
		default:

			coins, err := b.market.GetCoins(ctx)
			if err != nil {
				b.log.Errorf("Failed to fetch coins: %s.", err)
				continue
			}

			orders := b.db.GetOrders(models.BuyOrder, b.market.Name())
			for _, boughtCoin := range orders {
				takeProfit := boughtCoin.Price + (boughtCoin.Price*boughtCoin.TakeProfit)/100
				stopLoss := boughtCoin.Price + (boughtCoin.Price*(-1*math.Abs(boughtCoin.StopLoss)))/100
				lastPrice := coins[boughtCoin.Symbol].Price
				buyPrice := boughtCoin.Price
				priceChangePercentage := (lastPrice - buyPrice) / buyPrice * 100

				// Check that the price is above the take profit and readjust SL and TP accordingly if trialing stop loss is used.
				if b.config.TradingOptions.UseTrailingStopLoss && lastPrice >= takeProfit {
					boughtCoin.StopLoss = boughtCoin.TakeProfit - b.config.TradingOptions.TrailingStopLoss
					boughtCoin.TakeProfit = priceChangePercentage + b.config.TradingOptions.TrailingTakeProfit
					b.log.Infof("Price of %s reached more than the trading profit (TP). Adjusting stop loss (SL) to %f and trading profit (TP) to %f.", boughtCoin.Symbol, boughtCoin.StopLoss, boughtCoin.TakeProfit)
					b.db.SaveOrder(boughtCoin)
					continue
				}

				// Verify that the price is below the stop loss or above take profit and sell the boughtCoin.
				if lastPrice <= stopLoss || lastPrice >= takeProfit {
					var profitOrLossText string
					if priceChangePercentage >= 0 {
						profitOrLossText = "profit"
					} else {
						profitOrLossText = "loss"
					}

					estimatedProfitLoss := (lastPrice - buyPrice) * boughtCoin.Volume * (1 - (b.config.TradingOptions.TradingFee * 2))
					estimatedProfitLossWithFees := b.config.TradingOptions.Quantity * (priceChangePercentage - (b.config.TradingOptions.TradingFee * 2)) / 100
					msg := fmt.Sprintf(
						"Selling %f %s. Estimated %s: $%s %s%% (w/ fees: $%s %s%%)",
						boughtCoin.Volume,
						boughtCoin.Symbol,
						profitOrLossText,
						strconv.FormatFloat(estimatedProfitLoss, 'f', 2, 64),
						strconv.FormatFloat(priceChangePercentage, 'f', 2, 64),
						strconv.FormatFloat(estimatedProfitLossWithFees, 'f', 2, 64),
						strconv.FormatFloat(priceChangePercentage-(b.config.TradingOptions.TradingFee*2), 'f', 2, 64),
					)

					b.log.Infow(
						msg,
						"symbol", boughtCoin.Symbol,
						"buyPrice", buyPrice,
						"currentPrice", lastPrice,
						"priceChangePercentage", priceChangePercentage,
						"tradingFee", b.config.TradingOptions.TradingFee*2,
						"quantity", b.config.TradingOptions.Quantity,
						"testMode", b.config.ScriptOptions.TestMode,
					)

					order := models.Order{
						Market:                b.market.Name(),
						Type:                  models.SellOrder,
						Volume:                boughtCoin.Volume,
						PriceChangePercentage: priceChangePercentage,
						EstimatedProfitLoss:   estimatedProfitLoss,
					}

					if b.config.ScriptOptions.TestMode {
						order.Order = market.Order{
							OrderID:         0,
							Symbol:          boughtCoin.Symbol,
							TransactionTime: time.Now(),
							Price:           lastPrice,
						}
						order.IsTestMode = true
					} else {
						sellOrder, err := b.market.Sell(ctx, boughtCoin.Symbol, boughtCoin.Volume)
						if err != nil {
							b.log.Errorf("Failed to sell %s: %s.", boughtCoin.Symbol, err)
							continue
						}

						order.Order = sellOrder
					}

					b.db.SaveOrder(order)
					b.db.DeleteOrder(boughtCoin)
				}
			}

			time.Sleep(time.Second * time.Duration(b.config.TradingOptions.SellTimeout))
		}
	}
}

// updateLatestCoins fetches the latest coins from the market and appends them to the volatilityWindow.
func (b *Bot) updateLatestCoins(ctx context.Context) error {
	b.log.Info("Fetching latest coins.")

	coins, err := b.market.GetCoins(ctx)
	if err != nil {
		return err
	}

	b.volatilityWindow.AddRecord(coins)

	return nil
}

// convertVolume converts the volume given in the configured quantity from base currency (USDT) to each coin's volume.
func (b *Bot) convertVolume(ctx context.Context, quantity float64, volatileCoin market.VolatileCoin) (float64, error) {
	var stepSize float64

	// Get the step size of the coin from the local cache if it exists or from Binance if it doesn't (yet).
	// The step size never changes, so it's safe to cache it forever.
	// This approach avoids an additional API request to Binance per trade.
	cache, ok := b.db.GetCache(volatileCoin.Symbol)
	if ok {
		stepSize = cache.StepSize
	} else {
		info, err := b.market.GetSymbolInfo(ctx, volatileCoin.Symbol)
		if err != nil {
			return 0, err
		}

		stepSize = info.StepSize

		b.db.SaveCache(models.Cache{Symbol: volatileCoin.Symbol, StepSize: stepSize})
	}

	volume := quantity / volatileCoin.Price

	// Round the volume to the step size of the coin.
	if stepSize != 0 {
		volume = utils.RoundStepSize(volume, stepSize)
	}

	return volume, nil
}
