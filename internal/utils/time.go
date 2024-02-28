package utils

import (
	"time"
)

// CalculateTimeDuration calculates the time duration for the given minutes and interval.
//
// For example, if minutes is 60 and interval is 5, the result will be 12 minutes because the division of `minutes` (60) by `interval` (5) is 12.
// The result (12) is then multiplied by time.Minute, which represents 1 minute in Go's time package.
func CalculateTimeDuration(minutes, interval int) time.Duration {
	return time.Duration(float64(minutes) / float64(interval) * float64(time.Minute))
}
