package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResidual(t *testing.T) {
	// given
	preds := []int{1, 2, 3, 5, 10, 7}
	y____ := []int{1, 3, 4, 10, 2, 6}

	// when
	min, max := residual(preds, y____)

	// then
	assert.Equal(t, -8, min)
	assert.Equal(t, 5, max)
}

func TestScale(t *testing.T) {
	// given

	// when
	scaled := scale(0.5, 10)
	assert.Equal(t, 4, scaled)
	// when
	scaled = scale(0.49398382, 7)
	assert.Equal(t, 2, scaled)
	// when
	scaled = scale(0.1, 11)
	assert.Equal(t, 0, scaled)
}
