package index

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/BenJoyenConseil/rmi/estimate"
	"github.com/BenJoyenConseil/rmi/estimate/linear"
	"github.com/BenJoyenConseil/rmi/estimate/polynomial"
	"github.com/BenJoyenConseil/rmi/search"
)

/*
LearnedIndex is an index structure that use inference to locate keys
*/
type LearnedIndex struct {
	M                        estimate.Estimator
	ST                       *search.SortedTable
	Len                      int
	MinErrBound, MaxErrBound int
}

/*
New return an LearnedIndex fitted over the dataset with a linear regression algorythm
*/
func New(dataset []float64) *LearnedIndex {

	st := search.NewSortedTable(dataset)
	// store.Flush(st)

	x, y := linear.Cdf(st.Keys)
	len_ := len(dataset)
	m := linear.Fit(x, y)
	guesses := make([]int, len_)
	scaledY := make([]int, len_)
	maxErr, minErr := 0, 0
	for i, k := range x {
		guesses[i] = scale(m.Predict(k), len_)
		scaledY[i] = scale(y[i], len_)
		residual := residual(guesses[i], scaledY[i])
		if residual > maxErr {
			maxErr = residual
		} else if residual < minErr {
			minErr = residual
		}
	}
	return &LearnedIndex{M: m, Len: len_, ST: st, MinErrBound: minErr, MaxErrBound: maxErr}
}

func NewCubic(dataset []float64) *LearnedIndex {

	st := search.NewSortedTable(dataset)

	x, y := polynomial.Cdf(st.Keys)
	dropDupX, dropDupY := polynomial.DropDuplicates(x, y)
	dropDupX, dropDupY = polynomial.Segment(dropDupX, dropDupY, 3)

	len_ := len(dataset)
	m := polynomial.Fit(dropDupX, dropDupY)
	guesses := make([]int, len_)
	scaledY := make([]int, len_)
	maxErr, minErr := 0, 0
	for i, k := range x {
		guesses[i] = scale(m.Predict(k), len_)
		scaledY[i] = scale(y[i], len_)
		residual := scaledY[i] - guesses[i]
		if residual > maxErr {
			maxErr = residual
		} else if residual < minErr {
			minErr = residual
		}
	}

	w := tabwriter.NewWriter(os.Stdout, 8, 0, 1, ' ', tabwriter.AlignRight)
	fmt.Fprintln(w, "i\tx\ty\tpred\tresidual\t")
	for i := 0; i < len(x); i++ {
		fmt.Fprintf(w, "%d", i+1)
		fmt.Fprintf(w, "\t%f", x[i])
		fmt.Fprintf(w, "\t%d", scale(y[i], len_))
		fmt.Fprintf(w, "\t%d", scale(m.Predict(x[i]), len_))
		fmt.Fprintf(w, "\t%d", scale(y[i], len_)-scale(m.Predict(x[i]), len_))
		fmt.Fprintln(w, "\t")
	}
	fmt.Fprintln(w)
	w.Flush()

	return &LearnedIndex{M: m, Len: len_, ST: st, MinErrBound: minErr, MaxErrBound: maxErr}
}

/*
GuessIndex return the predicted position of the key in the index
and upper / lower positions' search interval. Guess, lower and upper
always have values between 0 and len(keys)-1
*/
func (idx *LearnedIndex) GuessIndex(key float64) (guess, lower, upper int) {
	guess = scale(idx.M.Predict(key), idx.Len)
	lower = idx.MinErrBound + guess
	if lower < 0 {
		lower = 0
	} else if lower > idx.Len-1 {
		lower = idx.Len - 1
	}
	upper = guess + idx.MaxErrBound
	if upper > idx.Len-1 {
		upper = idx.Len - 1
	} else if upper < 0 {
		upper = 0
	}

	if guess < 0 {
		guess = 0
	} else if guess > idx.Len-1 {
		guess = idx.Len - 1
	}
	return guess, lower, upper
}

/*
Lookup return the first offsets of the key or err if the key is not found in the index
*/
func (idx *LearnedIndex) Lookup(key float64) (offsets []int, err error) {
	guess, lower, upper := idx.GuessIndex(key)
	i := 0
	// k, o, err := store.Get(guess_i)
	// st, err := store.STExtract(guess_i+1, upper+1)

	if key > idx.ST.Keys[guess] {
		subKeys := idx.ST.Keys[guess+1 : upper+1]
		i = sort.SearchFloat64s(subKeys, key) + guess + 1
	} else if key <= idx.ST.Keys[guess] {
		subKeys := idx.ST.Keys[lower : guess+1]
		i = sort.SearchFloat64s(subKeys, key) + lower
	}

	// iterate to get all equal keys
	for ; i < upper+1; i++ {
		if idx.ST.Keys[i] == key {
			offsets = append(offsets, idx.ST.Offsets[i])
		} else {
			break
		}
	}

	if len(offsets) == 0 {
		err = fmt.Errorf("The following key <%f> is not found in the index", key)
	}

	return offsets, err
}
