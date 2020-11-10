package linear

import (
	"math"
)

type LearnedIndex struct {
	m        *RegressionModel
	len      int
	maxError int
}

/*
New return an Index fitted and data
*/
func New(dataset []float64) *LearnedIndex {
	x, y := cdf(dataset)
	len := len(dataset)
	m := Fit(x, y)
	maxErr := 0
	for i, k := range x {
		guessPosition := math.Round(m.Predict(k)*float64(len) - 1)
		truePosition := y[i]*float64(len) - 1
		residual := math.Sqrt(math.Pow(truePosition-guessPosition, 2))
		if float64(maxErr) < residual {
			maxErr = int(residual)
		}
	}
	return &LearnedIndex{m: m, len: len, maxError: maxErr}
}

/*
GuessIndex return the predicted position of the key in the index
and upper / lower positions' search interval
*/
func (idx *LearnedIndex) GuessIndex(key float64) (guess, lower, upper int) {
	cdfPred := idx.m.Predict(key)
	scaleCDF := cdfPred*float64(idx.len) - 1
	roundPosition := math.Round(scaleCDF)
	guess = int(roundPosition)
	lower = guess - idx.maxError
	upper = guess + idx.maxError
	return guess, lower, upper
}

func (idx *LearnedIndex) Lookup(key float64) (position int) {
	return -1
}
