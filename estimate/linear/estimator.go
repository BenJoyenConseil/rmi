package linear

import (
	"math"

	"gonum.org/v1/gonum/stat"
)

/*
Model is a LinearRegression model
*/
type RegressionModel struct {
	Intercept, Slope float64
}

/*
Fit return a Model structure fitted with alpha and beta
of a LinearRegression applied on x. Y is the CDF value
*/
func Fit(x, y []float64) *RegressionModel {
	alpha, beta := stat.LinearRegression(x, y, nil, false)
	if math.IsNaN(alpha) || math.IsNaN(beta) {
		alpha = stat.Mean(y, nil)
		beta = 0
	}
	return &RegressionModel{Intercept: alpha, Slope: beta}
}

/*
Predict the CDF result of a given x
*/
func (m *RegressionModel) Predict(x float64) (predCDF float64) {
	y := m.Intercept + m.Slope*x
	return y
}

/*
Return the x array sorted and the y array containing
empirical CDF value foreach x's value. len(x)=len(y)
*/
func Cdf(x []float64) (sortedX, y []float64) {
	for _, i := range x {
		yi := stat.CDF(i, stat.Empirical, x, nil)
		y = append(y, yi)
	}
	return x, y
}
