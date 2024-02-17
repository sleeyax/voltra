package utils

import (
	"time"
)

// CalculateTimeDelta calculates the interval between two time points.
func CalculateTimeDelta(minutes, interval int) time.Duration {
	return time.Duration(float64(minutes) / float64(interval) * float64(time.Minute))
}
