package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRoundStepSize(t *testing.T) {
	assert.Equal(t, 1.1, RoundStepSize(1.1, 0.01))
	assert.Equal(t, 0.2, RoundStepSize(0.2, 0.01))
	assert.Equal(t, 26.0, RoundStepSize(25.9, 1.0))
}
