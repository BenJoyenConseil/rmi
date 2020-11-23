package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResidual(t *testing.T) {

	// when
	r := residual(1, 1)
	// then
	assert.Equal(t, 0, r)
	// when
	r = residual(2, 3)
	// then
	assert.Equal(t, 1, r)
	// when
	r = residual(10, 2)
	// then
	assert.Equal(t, -8, r)
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
