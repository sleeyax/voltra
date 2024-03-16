package main

import (
	"context"
	"fmt"
	"github.com/sleeyax/voltra/internal/bot"
	"github.com/sleeyax/voltra/internal/config"
	"github.com/sleeyax/voltra/internal/database"
	"github.com/sleeyax/voltra/internal/market"
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

	b := bot.New(&c, market.NewBinance(c), database.NewSqliteDatabase("voltra.db", c.LoggingOptions))
	b.Start(ctx)
}
