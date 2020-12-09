package index

import (
	"fmt"
	"testing"

	"github.com/BenJoyenConseil/rmi/linear"
	"github.com/BenJoyenConseil/rmi/search"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	idx := New(keys)

	// then
	assert.Equal(t, 7, idx.Len)
	assert.IsType(t, &linear.RegressionModel{}, idx.M)
	assert.Equal(t, 2, idx.MaxErrBound)
	assert.Equal(t, -2, idx.MinErrBound)
	// check the order
	assert.Equal(t, []float64{2.5, 2.98, 3, 3, 3.14, 5, 10}, idx.ST.Keys)
	assert.Equal(t, []int{5, 6, 1, 2, 3, 0, 4}, idx.ST.Offsets)
}

func TestGuessIndex(t *testing.T) {
	idx := &LearnedIndex{
		M:           &linear.RegressionModel{Intercept: .23119036646681634, Slope: .08523040437506509},
		Len:         7,
		MaxErrBound: 2,
		MinErrBound: -3,
	}

	// when
	guess, lower, upper := idx.GuessIndex(5)
	assert.Equal(t, 1, lower)
	assert.Equal(t, 4, guess)
	assert.Equal(t, 6, upper)

	// when upper < 0 ==> 0
	guess, lower, upper = idx.GuessIndex(-10)
	assert.Equal(t, 0, lower)
	assert.Equal(t, 0, guess)
	assert.Equal(t, 0, upper)

	// when guess > len-1 ==> 6
	guess, lower, upper = idx.GuessIndex(10.)
	assert.Equal(t, 4, lower)
	assert.Equal(t, 6, guess)
	assert.Equal(t, 6, upper)
}

func TestLookup(t *testing.T) {
	// given
	idx := &LearnedIndex{
		M:           &linear.RegressionModel{Intercept: .23119036646681634, Slope: .08523040437506509},
		Len:         7,
		MaxErrBound: 2,
		MinErrBound: -2,
		//	keys : {5, 3, 3, 3.14, 10, 2.5, 2.98}
		//  sort : {2.5, 2.98, 3, 3, 3.14, 5, 10}
		//  posi : {0,   1,    2, 3, 4,    5, 6}
		ST: &search.SortedTable{
			Keys:    []float64{2.5, 2.98, 3, 3, 3.14, 5, 10},
			Offsets: []int{5, 6, 1, 2, 3, 0, 4},
		},
	}

	// when
	offsets, err := idx.Lookup(2.5)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{5}, offsets)
	// when
	offsets, err = idx.Lookup(2.98)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{6}, offsets)
	// when
	offsets, err = idx.Lookup(3.)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2}, offsets)
	// when
	offsets, err = idx.Lookup(3.14)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{3}, offsets)
	// when
	offsets, err = idx.Lookup(5.)
	assert.Nil(t, err)
	assert.Equal(t, []int{0}, offsets)
	// when
	offsets, err = idx.Lookup(10.)
	assert.Nil(t, err)
	assert.Equal(t, []int{4}, offsets)

	// when not in the index
	offsets, err = idx.Lookup(199.)
	assert.NotNil(t, err)
	assert.Nil(t, offsets)
}

func ExampleLearnedIndex() {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	index := New(keys)

	for _, k := range keys {
		offset, err := index.Lookup(k)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf("The key %f is located %d\n", k, offset)
	}

	// Output:
	// The key 2.500000 is located [5]
	// The key 2.980000 is located [6]
	// The key 3.000000 is located [1 2]
	// The key 3.000000 is located [1 2]
	// The key 3.140000 is located [3]
	// The key 5.000000 is located [0]
	// The key 10.000000 is located [4]
}
