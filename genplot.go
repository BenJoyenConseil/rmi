package main

import (
	"image/color"
	"log"
	"math"
	"rmi/index"
	"rmi/linear"
	"testing"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func TestExploration(t *testing.T) {

	// load the age column and parse values into float64 values
	ageColumn := extractAgeColumn("data/people.csv")
	log.Println("Values to index:", ageColumn)

	// create an index over the age column
	index := index.New(ageColumn)

	linearRegFn := func(x float64) float64 { return index.M.Predict(x)*float64(index.Len) - 1 }
	x, y := linear.Cdf(ageColumn)
	idxFromCDF := func(i float64) float64 { return stat.CDF(i, stat.Empirical, x, nil)*float64(index.Len) - 1 }

	p, _ := plot.New()
	p.Title.Text = "Learned Index RMI"
	p.X.Label.Text = "Age"
	p.Y.Label.Text = "Index"

	courbeKeys := plotter.XYs{}
	courbePreds := plotter.XYs{}
	for i, k := range x {
		courbeKeys = append(courbeKeys, plotter.XY{X: k, Y: float64(i)})
		pred := math.Round(index.M.Predict(k)*float64(index.Len) - 1)
		yIdx := y[i]*float64(index.Len) - 1
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
	p.X.Max = 100
	p.Y.Min = 0
	p.Y.Max = 10
	plotutil.AddScatters(p, s, "preds")
	p.Save(4*vg.Inch, 4*vg.Inch, "assets/plot.svg")
}
