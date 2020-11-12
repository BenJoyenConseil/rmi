package index

import (
	"rmi/linear"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	idx := New(keys)

	// then
	assert.Equal(t, 7, idx.Len)
	assert.Equal(t, .23119036646681634, idx.M.Intercept)
	assert.Equal(t, .08523040437506509, idx.M.Slope)
	assert.Equal(t, 2, idx.maxError)
	// check the order
	assert.Equal(t, &kv{key: 2.5, offset: 5}, idx.sortedTable[0])
	assert.Equal(t, &kv{key: 2.98, offset: 6}, idx.sortedTable[1])
	assert.Equal(t, &kv{key: 3, offset: 1}, idx.sortedTable[2])
	assert.Equal(t, &kv{key: 3, offset: 2}, idx.sortedTable[3])
	assert.Equal(t, &kv{key: 3.14, offset: 3}, idx.sortedTable[4])
	assert.Equal(t, &kv{key: 5, offset: 0}, idx.sortedTable[5])
	assert.Equal(t, &kv{key: 10, offset: 4}, idx.sortedTable[6])
}

func TestGuessIndex(t *testing.T) {
	idx := &LearnedIndex{
		M:        &linear.RegressionModel{Intercept: .23119036646681634, Slope: .08523040437506509},
		Len:      7,
		maxError: 2,
	}

	// when
	guess, lower, upper := idx.GuessIndex(2.95)
	assert.Equal(t, 2, guess)
	assert.Equal(t, 0, lower)
	assert.Equal(t, 4, upper)

	// when guess < 0 ==> -1
	guess, lower, upper = idx.GuessIndex(-2.95)
	assert.Equal(t, 0, guess)
	assert.Equal(t, 0, lower)
	assert.Equal(t, 2, upper)

	// when guess > len-1 ==> 7
	guess, lower, upper = idx.GuessIndex(10.)
	assert.Equal(t, 6, guess)
	assert.Equal(t, 4, lower)
	assert.Equal(t, 6, upper)
}

func TestLookup(t *testing.T) {
	// given
	idx := &LearnedIndex{
		M:        &linear.RegressionModel{Intercept: .23119036646681634, Slope: .08523040437506509},
		Len:      7,
		maxError: 2,
		//	keys : {5, 3, 3, 3.14, 10, 2.5, 2.98}
		//  sort : {2.5, 2.98, 3, 3, 3.14, 5, 10}
		//  posi : {0,   1,    2, 3, 4,    5, 6}
		sortedTable: []*kv{
			&kv{key: 2.5, offset: 5},
			&kv{key: 2.98, offset: 6},
			&kv{key: 3., offset: 1},
			&kv{key: 3., offset: 2},
			&kv{key: 3.14, offset: 3},
			&kv{key: 5., offset: 0},
			&kv{key: 10., offset: 4},
		},
	}

	// when
	offset, err := idx.Lookup(2.98)
	assert.Nil(t, err)
	assert.Equal(t, 6, offset)

	// when
	offset, err = idx.Lookup(3.)
	assert.Nil(t, err)
	assert.Equal(t, 1, offset)

	// when
	offset, err = idx.Lookup(3.14)
	assert.Nil(t, err)
	assert.Equal(t, 3, offset)

	// when
	offset, err = idx.Lookup(10.)
	assert.Nil(t, err)
	assert.Equal(t, 4, offset)

	// when not in the index
	offset, err = idx.Lookup(199.)
	assert.NotNil(t, err)
	assert.Equal(t, -1, offset)
}

func TestBinarySearch(t *testing.T) {
	// given
	kspace := []*kv{
		&kv{key: 1, offset: 0},
		&kv{key: 1, offset: 1},
		&kv{key: 1, offset: 2},
		&kv{key: 3, offset: 3},
		&kv{key: 3.14, offset: 4},
		&kv{key: 5, offset: 123},
		&kv{key: 5, offset: 124},
		&kv{key: 7.544, offset: 125},
		&kv{key: 7.544, offset: 126},
		&kv{key: 7.599, offset: 127},
	}

	// when
	offset, err := binarySearch(1, kspace)
	assert.Contains(t, []int{0, 1, 2}, offset)
	assert.Nil(t, err)

	// when
	offset, err = binarySearch(3, kspace)
	assert.Equal(t, 3, offset)
	assert.Nil(t, err)

	// when
	offset, err = binarySearch(3.14, kspace)
	assert.Equal(t, 4, offset)
	assert.Nil(t, err)

	// when
	offset, err = binarySearch(5, kspace)
	assert.Contains(t, []int{123, 124}, offset)
	assert.Nil(t, err)

	// when
	offset, err = binarySearch(7.544, kspace)
	assert.Equal(t, 125, offset)
	assert.Nil(t, err)

	// when
	offset, err = binarySearch(7.599, kspace)
	assert.Equal(t, 127, offset)
	assert.Nil(t, err)

	// when not found
	offset, err = binarySearch(14, kspace)
	assert.NotNil(t, err)
}
