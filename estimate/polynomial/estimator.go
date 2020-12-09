package polynomial

import (
	"sort"

	"gonum.org/v1/gonum/interp"
	"gonum.org/v1/gonum/stat"
)

func Fit(x, y []float64) interp.FittablePredictor {
	pl := &interp.PiecewiseLinear{}
	m := interp.FittablePredictor(pl)
	m.Fit(x, y)
	return m
}

/*
Return the x array sorted and the y array containing
empirical CDF value foreach x's value. len(x)=len(y)
*/
func Cdf(xs []float64) (sortedX, y []float64) {
	mapX := map[float64]interface{}{}
	if !sort.Float64sAreSorted(xs) {
		sort.Float64s(xs)
	}
	for i, x := range xs {
		if _, ok := mapX[x]; ok {
			continue
		}
		prev := 0
		if i > 0 {
			prev = i
		}
		yi := stat.CDF(xs[prev], stat.Empirical, xs, nil)
		y = append(y, yi)
	}
	return xs, y
}

func Segment(xs, ys []float64, n int) (segX, segY []float64) {
	size := len(xs) / n
	for i := 0; i < n; i++ {
		segX = append(segX, xs[i*size])
		segY = append(segY, ys[i*size])
	}
	segX = append(segX, xs[len(xs)-1])
	segY = append(segY, ys[len(ys)-1])
	return segX, segY
}

func DropDuplicates(sortedX, sortedY []float64) (sortedSetX, sortedSetY []float64) {
	mapX := map[float64]interface{}{}

	for i := range sortedX {
		if _, ok := mapX[sortedX[i]]; ok {
			continue
		}
		mapX[sortedX[i]] = nil
		sortedSetX = append(sortedSetX, sortedX[i])
		sortedSetY = append(sortedSetY, sortedY[i])
	}
	return sortedSetX, sortedSetY
}
