package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSortedTable(t *testing.T) {
	// given
	unsortedkeys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	st := NewSortedTable(unsortedkeys)

	// then
	assert.Equal(t, []float64{2.5, 2.98, 3, 3, 3.14, 5, 10}, st.Keys)
	assert.Equal(t, []int{5, 6, 1, 2, 3, 0, 4}, st.Offsets)
}
