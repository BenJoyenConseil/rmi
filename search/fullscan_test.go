package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFullScanLookup(t *testing.T) {
	// given
	st := &SortedTable{
		Keys:    []float64{.2342, 1.234, 2., 2., 3., 3., 10., 28},
		Offsets: []int{3, 2, 1, 4, 6, 0, 7, 5},
	}

	// when
	o, err := FullScanLookup(3., st)
	// then
	assert.ElementsMatch(t, o, []int{6, 0})
	assert.NoError(t, err)

	// when
	o, err = FullScanLookup(.2342, st)
	// then
	assert.ElementsMatch(t, o, []int{3})
	assert.NoError(t, err)

	// when
	o, err = FullScanLookup(2., st)
	// then
	assert.ElementsMatch(t, o, []int{1, 4})
	assert.NoError(t, err)

	// when
	o, err = FullScanLookup(10., st)
	// then
	assert.ElementsMatch(t, o, []int{7})
	assert.NoError(t, err)

	// when
	o, err = FullScanLookup(4., st)
	// then
	assert.Nil(t, o)
	assert.Error(t, err)
}
