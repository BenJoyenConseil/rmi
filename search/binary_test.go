package search

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinarySearch(t *testing.T) {
	// given
	st := &SortedTable{
		Keys:    []float64{.2342, 1.234, 2., 2., 3., 3., 10., 28},
		Offsets: []int{3, 2, 1, 4, 6, 0, 7, 5},
	}

	// when
	o, err := BinarySearchLookup(3., st)
	// then
	assert.ElementsMatch(t, o, []int{6, 0})
	assert.NoError(t, err)

	// when
	o, err = BinarySearchLookup(.2342, st)
	// then
	assert.ElementsMatch(t, o, []int{3})
	assert.NoError(t, err)

	// when
	o, err = BinarySearchLookup(2., st)
	// then
	assert.ElementsMatch(t, o, []int{1, 4})
	assert.NoError(t, err)

	// when
	o, err = BinarySearchLookup(10., st)
	// then
	assert.ElementsMatch(t, o, []int{7})
	assert.NoError(t, err)

	// when
	o, err = BinarySearchLookup(4., st)
	// then
	assert.Nil(t, o)
	assert.Error(t, err)
}

func TestInterpolationSearch(t *testing.T) {
	// given
	st := table{1, 2, 3, 4, 5, 5, 5, 6, 7}
	key := 5.
	middle := 6

	// when
	i := InterpolationSearch(key, st, middle)

	// then leftmost 5
	assert.Equal(t, 4, i)
}

func TestSearchTable(t *testing.T) {
	// given
	st := table{1, 2, 3, 4, 5, 5, 5, 6, 7}

	// when
	i := SearchTable(st, 5)

	// then leftmost 5
	assert.Equal(t, 4, i)
}

type table []float64

func (t table) Get(i int) (float64, int) { return t[i], -1 }
func (t table) Len() int                 { return len(t) }

func BenchmarkBinarySeach(b *testing.B) {
	keys := make([]float64, 10000)
	for i := 0; i < 10000; i++ {
		keys[i] = rand.Float64()
	}
	st := NewSortedTable(keys)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		SearchTable(st, rand.Float64())
	}
}
