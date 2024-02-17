package main

import (
	"context"
	"fmt"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/config"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
)

func main() {
	c, err := config.Load()
	if err != nil {
		panic(fmt.Errorf("failed to load config file: %w", err))
	}

	m := market.NewBinance(c)
	coins, err := m.GetCoins(context.Background())
	if err != nil {
		panic(fmt.Errorf("failed to get coins from market: %w", err))
	}

	fmt.Println(coins)
}
