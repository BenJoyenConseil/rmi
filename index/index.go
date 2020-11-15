package index

import (
	"fmt"
	"math"
	"rmi/linear"
	"sort"
)

type kv struct {
	key    float64
	offset int
}

type sortedTable struct {
	keys    []float64
	offsets []int
}

// Implement the Sort interface
type byKeys struct{ *sortedTable }

func (st byKeys) Len() int { return len(st.keys) }
func (st byKeys) Swap(i, j int) {
	st.keys[i], st.keys[j] = st.keys[j], st.keys[i]
	st.offsets[i], st.offsets[j] = st.offsets[j], st.offsets[i]
}
func (st byKeys) Less(i, j int) bool { return st.keys[i] < st.keys[j] }

func newSortedTable(x []float64) *sortedTable {
	keys, offsets := x, make([]int, len(x))
	for i := range x {
		offsets[i] = i
	}
	st := &sortedTable{keys: keys, offsets: offsets}
	sort.Sort(byKeys{st})
	return st
}

/*
LearnedIndex is an index structure that use inference to locate keys
*/
type LearnedIndex struct {
	M           *linear.RegressionModel
	ST          *sortedTable
	sortedTable []*kv
	Len         int
	maxError    int
}

/*
New return an LearnedIndex fitted over the dataset with a linear regression algorythm
*/
func New(dataset []float64) *LearnedIndex {

	st := newSortedTable(dataset)

	x, y := linear.Cdf(st.keys)
	len := len(dataset)
	m := linear.Fit(x, y)
	maxErr := 0
	for i, k := range x {
		guessPosition := math.Round(scale(m.Predict(k), len))
		truePosition := math.Round(scale(y[i], len))
		residual := math.Sqrt(math.Pow(truePosition-guessPosition, 2))
		if float64(maxErr) < residual {
			maxErr = int(residual)
		}

	}
	return &LearnedIndex{M: m, Len: len, maxError: maxErr, ST: st}
}

/*
scale return the CDF value x datasetLen -1 to get back the position in a sortedTable
*/
func scale(cdfVal float64, datasetLen int) float64 {
	return cdfVal*float64(datasetLen) - 1
}

/*
GuessIndex return the predicted position of the key in the index
and upper / lower positions' search interval. Guess, lower and upper
always have values between 0 and len(keys)-1
*/
func (idx *LearnedIndex) GuessIndex(key float64) (guess, lower, upper int) {
	guess = int(math.Round(scale(idx.M.Predict(key), idx.Len)))
	if guess < 0 {
		guess = 0
	} else if guess > idx.Len-1 {
		guess = idx.Len - 1
	}
	lower = guess - idx.maxError
	if lower < 0 {
		lower = 0
	}
	upper = guess + idx.maxError
	if upper > idx.Len-1 {
		upper = idx.Len - 1
	}
	return guess, lower, upper
}

/*
Lookup return the first offsets of the key or err if the key is not found in the index
*/
func (idx *LearnedIndex) Lookup(key float64) (offsets []int, err error) {
	guess, lower, upper := idx.GuessIndex(key)
	i := 0

	if idx.ST.keys[guess] == key {
		i = guess
	} else if idx.ST.keys[guess] < key {
		subKeys := idx.ST.keys[guess+1 : upper+1]
		i = sort.SearchFloat64s(subKeys, key) + guess + 1
	} else {
		subKeys := idx.ST.keys[lower:guess]
		i = sort.SearchFloat64s(subKeys, key) + lower
	}

	// iterate to get all equal keys
	for i := i; i < len(idx.ST.keys); i++ {
		if idx.ST.keys[i] == key {
			offsets = append(offsets, idx.ST.offsets[i])
		} else {
			break
		}
	}

	if len(offsets) == 0 {
		err = fmt.Errorf("The following key <%f> is not found in the index", key)
	}

	return offsets, err
}
