package main

import (
	"context"
	"fmt"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/bot"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"os"
	"os/signal"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	c, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %w", err))
	}

	b := bot.New(&c, market.NewBinance(c), database.NewSqliteDatabase("data.db"))
	b.Start(ctx)
}
