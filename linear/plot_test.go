package linear

import (
	"image/color"
	"log"
	"math"
	"testing"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func TestExploration(t *testing.T) {
	keys := []float64{2.5, 2.98, 3, 3, 3.14, 5, 10}
	lenKeys := float64(len(keys))
	m := &RegressionModel{Intercept: 0.23119036646681634, Slope: 0.08523040437506509}
	linearRegFn := func(x float64) float64 { return m.Predict(x)*lenKeys - 1 }
	x, y := Cdf(keys)
	idxFromCDF := func(i float64) float64 { return stat.CDF(i, stat.Empirical, x, nil)*lenKeys - 1 }

	p, _ := plot.New()
	p.Title.Text = "Learned Index RMI"
	p.X.Label.Text = "Keys"
	p.Y.Label.Text = "Index"

	courbeKeys := plotter.XYs{}
	maxErr, minErr := 0., 0.
	for i, k := range x {
		courbeKeys = append(courbeKeys, plotter.XY{X: k, Y: float64(i)})
		pred := math.Round(m.Predict(k)*lenKeys - 1)
		yIdx := y[i]*lenKeys - 1
		residual := yIdx - pred
		if residual > maxErr {
			maxErr = residual
		} else if residual < minErr {
			minErr = residual
		}
		log.Println("y", yIdx, "guess", pred, "err:", residual)
	}
	approxFn := plotter.NewFunction(linearRegFn)
	approxFn.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	approxFn.Width = vg.Points(2)
	approxFn.Color = color.RGBA{G: 255, A: 255}

	cdfFn := plotter.NewFunction(idxFromCDF)
	cdfFn.Width = vg.Points(1)
	cdfFn.Color = color.RGBA{A: 255, B: 255}

	maxErrFn := plotter.NewFunction(func(x float64) float64 { return m.Predict(x)*lenKeys - 1 + maxErr })
	maxErrFn.Dashes = []vg.Length{vg.Points(1), vg.Points(5)}
	maxErrFn.Width = vg.Points(1)
	maxErrFn.Color = plotutil.DefaultColors[6]
	p.Add(maxErrFn)
	p.Legend.Add("max error", maxErrFn)
	minErrFn := plotter.NewFunction(func(x float64) float64 { return m.Predict(x)*lenKeys - 1 + minErr })
	minErrFn.Dashes = []vg.Length{vg.Points(1), vg.Points(5)}
	minErrFn.Width = vg.Points(1)
	minErrFn.Color = plotutil.DefaultColors[4]
	p.Add(minErrFn)
	p.Legend.Add("min error", minErrFn)

	// s, _ := plotter.NewScatter(courbePreds)
	// s.Color = color.RGBA{G: 255, A: 255}
	// s.Shape = draw.PyramidGlyph{}
	// plotutil.AddScatters(p, s, "preds")

	plotutil.AddLinePoints(p, "Keys", courbeKeys)
	p.Add(approxFn)
	p.Legend.Add("Approx (lr)", approxFn)
	p.Add(cdfFn)
	p.Legend.Add("CDF", cdfFn)
	p.X.Min = 0
	p.X.Max = 12
	p.Y.Min = 0
	p.Y.Max = 12
	p.Add(plotter.NewGrid())
	p.Save(4*vg.Inch, 4*vg.Inch, "plot.jpeg")
}
