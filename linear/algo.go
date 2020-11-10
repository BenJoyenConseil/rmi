package linear

import (
	"log"
	"sort"

	"gonum.org/v1/gonum/stat"
)

/*
Model is a LinearRegression model
*/
type RegressionModel struct {
	intercept, slope float64
}

/*
Fit return a Model structure fitted with alpha and beta
of a LinearRegression applied on x. Y is the CDF value
*/
func Fit(x, y []float64) *RegressionModel {
	alpha, beta := stat.LinearRegression(x, y, nil, false)
	return &RegressionModel{intercept: alpha, slope: beta}
}

/*
Predict the CDF result of a given x
*/
func (m *RegressionModel) Predict(x float64) (predCDF float64) {
	y := m.intercept + m.slope*x
	return y
}

/*
Return the x array sorted and the y array containing
empirical CDF value foreach x's value. len(x)=len(y)
*/
func cdf(x []float64) (sortedX, y []float64) {
	if !sort.Float64sAreSorted(x) {
		sort.Float64s(x)
	}
	for _, i := range x {
		yi := stat.CDF(i, stat.Empirical, x, nil)
		log.Println("CDF de ", i, "=", yi)
		y = append(y, yi)
	}
	return x, y
}
