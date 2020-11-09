package linear

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func TestFit(t *testing.T) {
	// given
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	m := Fit(keys)

	// then
	assert.NotNil(t, m)
	assert.Equal(t, m.intercept, 0.23119036646681634)
	assert.Equal(t, m.slope, 0.08523040437506509)
	assert.Equal(t, m.lenX, 7)
}

func TestPredict(t *testing.T) {
	// given
	key := 5.0
	m := &Model{intercept: 0.23119036646681634, slope: 0.08523040437506509, lenX: 7}

	// when
	predictedIdx := m.Predict(key)

	// then
	assert.Equal(t, 4, predictedIdx)
}

func TestCDF(t *testing.T) {
	// given
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	idx, y := cdf(keys)

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

func TestExploration(t *testing.T) {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}
	m := &Model{intercept: 0.23119036646681634, slope: 0.08523040437506509, lenX: 7}

	p, _ := plot.New()
	p.Title.Text = "Learned Index RMI"
	p.X.Label.Text = "Keys"
	p.Y.Label.Text = "Index"

	x, _ := cdf(keys)
	courbeKeys := plotter.XYs{}
	for i, k := range x {
		courbeKeys = append(courbeKeys, plotter.XY{X: k, Y: float64(i)})
	}
	approxFn := plotter.NewFunction(func(x float64) float64 { return (m.intercept+m.slope*x)*7 - 1 })
	approxFn.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	approxFn.Width = vg.Points(2)
	approxFn.Color = color.RGBA{G: 255, A: 255}

	cdfFn := plotter.NewFunction(func(i float64) float64 { return stat.CDF(i, stat.Empirical, x, nil)*7 - 1 })
	cdfFn.Width = vg.Points(1)
	cdfFn.Color = color.RGBA{A: 255, B: 255}

	plotutil.AddLinePoints(p, "Keys", courbeKeys)
	p.Add(approxFn)
	p.Legend.Add("Approx (lr)", approxFn)
	p.Add(cdfFn)
	p.Legend.Add("CDF", cdfFn)
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10
	p.Save(4*vg.Inch, 4*vg.Inch, "plot.jpeg")
}
