package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCalculateTimeDelta(t *testing.T) {
	var duration time.Duration

	duration = CalculateTimeDuration(2, 10)
	assert.Equal(t, 12.0, duration.Seconds())

	duration = CalculateTimeDuration(60, 5)
	assert.Equal(t, 12.0, duration.Minutes())
}
