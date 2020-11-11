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

type LearnedIndex struct {
	m           *linear.RegressionModel
	sortedTable []*kv
	len         int
	maxError    int
}

// ByKey implements sort.Interface for []*kv
type ByKey []*kv

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].key < a[j].key }

/*
New return an Index fitted and data
*/
func New(dataset []float64) *LearnedIndex {
	var sortedTable []*kv
	for i, row := range dataset {
		sortedTable = append(sortedTable, &kv{key: row, offset: i})
	}
	sort.Sort(ByKey(sortedTable))

	x, y := linear.Cdf(dataset)
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
	return &LearnedIndex{m: m, len: len, maxError: maxErr, sortedTable: sortedTable}
}

func scale(cdfVal float64, datasetLen int) float64 {
	return cdfVal*float64(datasetLen) - 1
}

/*
GuessIndex return the predicted position of the key in the index
and upper / lower positions' search interval
*/
func (idx *LearnedIndex) GuessIndex(key float64) (guess, lower, upper int) {
	guess = int(math.Round(scale(idx.m.Predict(key), idx.len)))
	if guess < 0 {
		guess = 0
	} else if guess > idx.len-1 {
		guess = idx.len - 1
	}
	lower = guess - idx.maxError
	if lower < 0 {
		lower = 0
	}
	upper = guess + idx.maxError
	if upper > idx.len-1 {
		upper = idx.len - 1
	}
	return guess, lower, upper
}

/*
Lookup return the actual value or err if the key is not found in the index
*/
func (idx *LearnedIndex) Lookup(key float64) (offset int, err error) {
	guess, lower, upper := idx.GuessIndex(key)

	if 0 <= guess && guess < idx.len {
		if idx.sortedTable[guess].key == key {
			return idx.sortedTable[guess].offset, nil
		} else if idx.sortedTable[guess].key < key {
			return binarySearch(key, idx.sortedTable[guess+1:upper+1])
		} else {
			return binarySearch(key, idx.sortedTable[lower:guess])
		}
	}

	return -1, fmt.Errorf("The following key <%f> is not found in the index", key)
}

/*
binarySearch implementation is for finding the leftmost element
*/
func binarySearch(key float64, searchSpace []*kv) (offeset int, err error) {
	L := 0
	R := len(searchSpace) - 1
	nIter := 0

	for L <= R {
		m := int(math.Floor(float64((L + R) / 2)))
		if searchSpace[m].key < key {
			L = m + 1
		} else if searchSpace[m].key > key {
			R = m - 1
		} else {
			return searchSpace[m].offset, nil
		}
		nIter++
	}
	return -1, fmt.Errorf("The following key <%f> is not found in the index", key)
}
