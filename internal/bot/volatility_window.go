package bot

import (
	"cmp"
	"github.com/sleeyax/voltra/internal/market"
	"math"
	"slices"
	"time"
)

// UnlimitedVolatilityWindowLength is a constant that can be used to indicate that the volatilityWindow should have an unlimited length.
// Warning: this should only be used for testing purposes. Growing the volatilityWindow indefinitely can lead to memory leaks.
const UnlimitedVolatilityWindowLength = 0

type VolatilityWindowRecord struct {
	time  time.Time
	coins market.CoinMap
}

type VolatilityWindow struct {
	records       []VolatilityWindowRecord
	volatileCoins market.VolatileCoins
	maxLength     int
}

// NewVolatilityWindow creates a new volatilityWindow of records.
// The volatilityWindow can be used to monitor the price changes of coins over time.
// It's a rolling window of records, so the volatilityWindow will never exceed the given max length.
func NewVolatilityWindow(maxLength int) *VolatilityWindow {
	return &VolatilityWindow{records: make([]VolatilityWindowRecord, 0), volatileCoins: make(market.VolatileCoins), maxLength: maxLength}
}

// Size returns the number of records in the volatilityWindow.
func (h *VolatilityWindow) Size() int {
	return len(h.records)
}

// AddRecord adds a new record to the volatilityWindow.
func (h *VolatilityWindow) AddRecord(coins market.CoinMap) {
	if l := h.Size(); l == h.maxLength && l != UnlimitedVolatilityWindowLength {
		// remove everything except the last record
		// h.records = h.records[h.maxLength-1:]
		h.records = make([]VolatilityWindowRecord, 0)
		h.volatileCoins = make(market.VolatileCoins)
	}
	h.records = append(h.records, VolatilityWindowRecord{time: time.Now(), coins: coins})
}

// GetLatestRecord returns the latest record in the volatilityWindow.
func (h *VolatilityWindow) GetLatestRecord() VolatilityWindowRecord {
	return h.records[len(h.records)-1]
}

func (h *VolatilityWindow) calculatePrice(symbol string, sign int, r1, r2 VolatilityWindowRecord) (float64, float64) {
	_, ok1 := r1.coins[symbol]
	_, ok2 := r2.coins[symbol]

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
		p1 = r1.coins[symbol].Price
		p2 = r2.coins[symbol].Price
	}

	return p1, p2
}

// Min returns the record with the lowest price for the given coin.
func (h *VolatilityWindow) Min(symbol string) VolatilityWindowRecord {
	return slices.MinFunc(h.records, func(r1, r2 VolatilityWindowRecord) int {
		p1, p2 := h.calculatePrice(symbol, 1, r1, r2)
		return cmp.Compare(p1, p2)
	})
}

// Max returns the record with the highest price for the given coin.
func (h *VolatilityWindow) Max(symbol string) VolatilityWindowRecord {
	return slices.MaxFunc(h.records, func(r1, r2 VolatilityWindowRecord) int {
		p1, p2 := h.calculatePrice(symbol, -1, r1, r2)
		return cmp.Compare(p1, p2)
	})
}

// IdentifyVolatileCoins returns the coins that have a price change of more than the given percentage.
// Returns a map of coin symbols and their respective price change percentage over the current time window of the volatilityWindow.
func (h *VolatilityWindow) IdentifyVolatileCoins(percentage float64) market.VolatileCoins {
	currentRecord := h.GetLatestRecord()

	for symbol, coin := range currentRecord.coins {
		minRecord := h.Min(symbol)
		maxRecord := h.Max(symbol)

		polarity := 1.0
		if minRecord.time.After(maxRecord.time) {
			polarity = -1.0
		}

		threshold := polarity * (maxRecord.coins[symbol].Price - minRecord.coins[symbol].Price) / minRecord.coins[symbol].Price * 100.0

		if threshold >= percentage {
			// only append the symbol if it's not already in the map
			_, ok := h.volatileCoins[symbol]
			if !ok {
				h.volatileCoins[symbol] = market.VolatileCoin{
					Coin:       coin,
					Percentage: threshold,
				}
			}
		}
	}

	return h.volatileCoins
}
