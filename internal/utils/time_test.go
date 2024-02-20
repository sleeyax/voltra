package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalculateTimeDelta(t *testing.T) {
	duration := CalculateTimeDelta(2, 10)
	assert.Equal(t, 12.0, duration.Seconds())
}
