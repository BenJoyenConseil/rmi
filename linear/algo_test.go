package linear

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFit(t *testing.T) {
	// given
	sortedX := []float64{2.5, 2.98, 3, 3, 3.14, 5, 10}
	y := []float64{0.14285714285714285, 0.2857142857142857, 0.5714285714285714, 0.5714285714285714, 0.7142857142857143, 0.8571428571428571, 1.0}
	// when
	m := Fit(sortedX, y)

	// then
	assert.NotNil(t, m)
	assert.Equal(t, m.Intercept, 0.23119036646681634)
	assert.Equal(t, m.Slope, 0.08523040437506509)
}

func TestFit_WhenAlphaBetaAreNaN_ShouldReturnMean(t *testing.T) {
	// given
	sortedX := []float64{3, 3, 3, 3}
	y := []float64{0.3, 0.4, 0.5, 0.6}
	// when
	m := Fit(sortedX, y)

	// then
	assert.NotNil(t, m)
	assert.Equal(t, 0.45, m.Intercept)
	assert.Equal(t, 0., m.Slope)
}

func TestPredict(t *testing.T) {
	// given
	m := &RegressionModel{Intercept: 0.23119036646681634, Slope: 0.08523040437506509}

	// when
	p2dot5 := m.Predict(2.5)
	p3 := m.Predict(3)
	p5 := m.Predict(5.)

	// then
	assert.Equal(t, .44426637740447905, p2dot5)
	assert.Equal(t, .4868815795920116, p3)
	assert.Equal(t, .6573423883421418, p5)
}

func TestCDF(t *testing.T) {
	// given
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	idx, y := Cdf(keys)

	// then
	assert.Equal(t, idx[0], 2.5)
	assert.Equal(t, idx[1], 2.98)
	assert.Equal(t, idx[2], 3.0)
	assert.Equal(t, idx[3], 3.0)
	assert.Equal(t, idx[4], 3.14)
	assert.Equal(t, idx[5], 5.0)
	assert.Equal(t, idx[6], 10.0)

	assert.Equal(t, y[0], 0.14285714285714285)
	assert.Equal(t, y[1], 0.2857142857142857)
	assert.Equal(t, y[2], 0.5714285714285714)
	assert.Equal(t, y[3], 0.5714285714285714)
	assert.Equal(t, y[4], 0.7142857142857143)
	assert.Equal(t, y[5], 0.8571428571428571)
	assert.Equal(t, y[6], 1.0)

	assert.Equal(t, (y[0]*float64(len(idx)))-1.0, 0.0)
}
