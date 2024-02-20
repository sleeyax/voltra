package bot

import (
	"cmp"
	"github.com/sleeyax/go-crypto-volatility-trading-bot/internal/market"
	"math"
	"slices"
	"time"
)

// UnlimitedHistoryLength is a constant that can be used to indicate that the history should have an unlimited length.
// Warning: this should only be used for testing purposes. Growing the history indefinitely can lead to memory leaks.
const UnlimitedHistoryLength = 0

type HistoryRecord struct {
	time  time.Time
	coins market.CoinMap
}

type VolatileCoins map[string]float64

type History struct {
	records       []HistoryRecord
	volatileCoins VolatileCoins
	maxLength     int
}

// NewHistory creates a new history of records.
// The history can be used to monitor the price changes of coins over time.
// It's a rolling window of records, so the history will never exceed the given max length.
func NewHistory(maxLength int) *History {
	return &History{records: make([]HistoryRecord, 0), volatileCoins: make(VolatileCoins), maxLength: maxLength}
}

// Size returns the number of records in the history.
func (h *History) Size() int {
	return len(h.records)
}

// AddRecord adds a new record to the history.
func (h *History) AddRecord(coins market.CoinMap) {
	if l := h.Size(); l == h.maxLength && l != UnlimitedHistoryLength {
		// remove everything except the last record
		// h.records = h.records[h.maxLength-1:]
		h.records = nil
		h.volatileCoins = nil
	}
	h.records = append(h.records, HistoryRecord{time: time.Now(), coins: coins})
}

// GetLatestRecord returns the latest record in the history.
func (h *History) GetLatestRecord() HistoryRecord {
	return h.records[len(h.records)-1]
}

func (h *History) calculatePrice(coinKey string, sign int, r1, r2 HistoryRecord) (float64, float64) {
	_, ok1 := r1.coins[coinKey]
	_, ok2 := r2.coins[coinKey]

	var p1 float64
	var p2 float64

	if !ok1 && !ok2 {
		p1 = math.Inf(sign)
		p2 = math.Inf(sign)
	} else if !ok1 {
		p1 = math.Inf(sign)
	} else if !ok2 {
		p2 = math.Inf(sign)
	} else {
		p1 = r1.coins[coinKey].Price
		p2 = r2.coins[coinKey].Price
	}

	return p1, p2
}

// Min returns the record with the lowest price for the given coin.
func (h *History) Min(coinKey string) HistoryRecord {
	return slices.MinFunc(h.records, func(r1, r2 HistoryRecord) int {
		p1, p2 := h.calculatePrice(coinKey, 1, r1, r2)
		return cmp.Compare(p1, p2)
	})
}

// Max returns the record with the highest price for the given coin.
func (h *History) Max(coinKey string) HistoryRecord {
	return slices.MaxFunc(h.records, func(r1, r2 HistoryRecord) int {
		p1, p2 := h.calculatePrice(coinKey, -1, r1, r2)
		return cmp.Compare(p1, p2)
	})
}

// IdentifyVolatileCoins returns the coins that have a price change of more than the given percentage.
// Returns a map of coin symbols and their respective price change percentage over the current time window of the history.
func (h *History) IdentifyVolatileCoins(percentage float64) VolatileCoins {
	currentRecord := h.GetLatestRecord()

	for coin := range currentRecord.coins {
		minRecord := h.Min(coin)
		maxRecord := h.Max(coin)

		polarity := 1.0
		if minRecord.time.After(maxRecord.time) {
			polarity = -1.0
		}

		threshold := polarity * (maxRecord.coins[coin].Price - minRecord.coins[coin].Price) / minRecord.coins[coin].Price * 100.0

		if threshold >= percentage {
			// only append the coin if it's not already in the map
			_, ok := h.volatileCoins[coin]
			if !ok {
				h.volatileCoins[coin] = threshold
			}
		}
	}

	return h.volatileCoins
}
