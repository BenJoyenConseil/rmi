package index

import (
	"image/color"

	"gonum.org/v1/gonum/floats"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

/*
Genplot takes an index and plots its keys, CDF, its approximation, and writes a plot.png file in assets/folder
*/
func Genplot(index *LearnedIndex, indexedCol []float64, plotfilepath string, roundedError bool) {
	linearRegFn := func(x float64) float64 { return index.M.Predict(x)*float64(index.Len) - 1 }
	idxFromCDF := func(i float64) float64 { return stat.CDF(i, stat.Empirical, index.ST.Keys, nil)*float64(index.Len) - 1 }

	p, _ := plot.New()
	p.Title.Text = "Learned Index RMI"
	p.X.Label.Text = "Age"
	p.Y.Label.Text = "Index"

	courbeKeys := plotter.XYs{}
	for i, k := range index.ST.Keys {
		courbeKeys = append(courbeKeys, plotter.XY{X: k, Y: float64(i)})
	}
	approxFn := plotter.NewFunction(linearRegFn)
	approxFn.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	approxFn.Width = vg.Points(2)
	approxFn.Color = color.RGBA{G: 255, A: 255}

	maxErrFn := plotter.NewFunction(func(x float64) float64 { _, _, upper := index.GuessIndex(x); return float64(upper) })
	if roundedError {
		maxErrFn = plotter.NewFunction(func(x float64) float64 { return float64(index.MaxErrBound) + index.M.Predict(x)*float64(index.Len) - 1 })
	}
	maxErrFn.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
	maxErrFn.Width = vg.Points(1)
	maxErrFn.Color = plotutil.SoftColors[2]
	p.Add(maxErrFn)
	p.Legend.Add("upper bound", maxErrFn)
	minErrFn := plotter.NewFunction(func(x float64) float64 { _, lower, _ := index.GuessIndex(x); return float64(lower) })
	if roundedError {
		minErrFn = plotter.NewFunction(func(x float64) float64 { return float64(index.MinErrBound) + index.M.Predict(x)*float64(index.Len) - 1 })
	}
	minErrFn.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
	minErrFn.Width = vg.Points(1)
	minErrFn.Color = plotutil.SoftColors[4]
	p.Add(minErrFn)
	p.Legend.Add("lower bound", minErrFn)

	cdfFn := plotter.NewFunction(idxFromCDF)
	cdfFn.Width = vg.Points(1)
	cdfFn.Color = color.RGBA{A: 255, B: 255}

	plotutil.AddLinePoints(p, "Keys", courbeKeys)
	p.Add(approxFn)
	p.Legend.Add("Approx (lr)", approxFn)
	p.Add(cdfFn)
	p.Legend.Add("CDF", cdfFn)
	p.X.Min = 0
	p.X.Max = floats.Max(index.ST.Keys)
	p.Y.Min = -float64(index.Len) / 10
	p.Y.Max = float64(index.Len) * 1.5
	p.Add(plotter.NewGrid())
	p.Save(4*vg.Inch, 4*vg.Inch, plotfilepath)
}
