package main

import (
	"context"
	"fmt"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/bot"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/database"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
)

func main() {
	c, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %w", err))
	}

	b := bot.New(&c, market.NewBinance(c), database.NewLocalDatabase())
	if err = b.Start(context.Background()); err != nil {
		panic(fmt.Errorf("failed to start bot: %w", err))
	}
}
