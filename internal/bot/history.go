package bot

import (
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"time"
)

type historyRecord struct {
	time  time.Time
	coins []market.Coin
}
