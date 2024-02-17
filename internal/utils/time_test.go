package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestCalculateTimeDelta(t *testing.T) {
	duration := CalculateTimeDelta(2, 10)
	assert.Equal(t, duration.Seconds(), 12.0)
}
