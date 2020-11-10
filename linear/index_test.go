package linear

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	idx := New(keys)

	// then
	assert.Equal(t, 7, idx.len)
	assert.Equal(t, .23119036646681634, idx.m.intercept)
	assert.Equal(t, .08523040437506509, idx.m.slope)
	assert.Equal(t, 2, idx.maxError)
}

func TestGuessIndex(t *testing.T) {
	idx := &LearnedIndex{
		m:        &RegressionModel{intercept: .23119036646681634, slope: .08523040437506509},
		len:      7,
		maxError: 2,
	}

	// when
	guess, lower, upper := idx.GuessIndex(2.95)

	// then
	assert.Equal(t, 2, guess)
	assert.Equal(t, 0, lower)
	assert.Equal(t, 4, upper)
}

func TestLookup(t *testing.T) {
	// when
	idx := &LearnedIndex{}
	
	// when
	idx.Lookup(1.)
}
