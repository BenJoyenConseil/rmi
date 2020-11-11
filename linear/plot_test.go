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
	"gonum.org/v1/plot/vg/draw"
)

func TestExploration(t *testing.T) {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}
	lenKeys := float64(len(keys))
	m := &RegressionModel{intercept: 0.23119036646681634, slope: 0.08523040437506509}
	linearRegFn := func(x float64) float64 { return m.Predict(x)*lenKeys - 1 }
	x, y := cdf(keys)
	idxFromCDF := func(i float64) float64 { return stat.CDF(i, stat.Empirical, x, nil)*lenKeys - 1 }

	p, _ := plot.New()
	p.Title.Text = "Learned Index RMI"
	p.X.Label.Text = "Keys"
	p.Y.Label.Text = "Index"

	courbeKeys := plotter.XYs{}
	courbePreds := plotter.XYs{}
	for i, k := range x {
		courbeKeys = append(courbeKeys, plotter.XY{X: k, Y: float64(i)})
		pred := math.Round(m.Predict(k)*lenKeys - 1)
		yIdx := y[i]*lenKeys - 1
		residual := math.Sqrt(math.Pow(yIdx-pred, 2.0))
		log.Println("y", yIdx, "guess", pred, "err:", residual)
		courbePreds = append(courbePreds, plotter.XY{X: k, Y: pred})
	}
	approxFn := plotter.NewFunction(linearRegFn)
	approxFn.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	approxFn.Width = vg.Points(2)
	approxFn.Color = color.RGBA{G: 255, A: 255}

	cdfFn := plotter.NewFunction(idxFromCDF)
	cdfFn.Width = vg.Points(1)
	cdfFn.Color = color.RGBA{A: 255, B: 255}

	s, _ := plotter.NewScatter(courbePreds)
	s.Color = color.RGBA{G: 255, A: 255}
	s.Shape = draw.PyramidGlyph{}

	plotutil.AddLinePoints(p, "Keys", courbeKeys)
	p.Add(approxFn)
	p.Legend.Add("Approx (lr)", approxFn)
	p.Add(cdfFn)
	p.Legend.Add("CDF", cdfFn)
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10
	plotutil.AddScatters(p, s, "preds")
	p.Save(4*vg.Inch, 4*vg.Inch, "plot.jpeg")
}
