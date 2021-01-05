package search

import (
	"fmt"
	"math"
	"sort"
)

func BinarySearchLookup(key float64, st *SortedTable) (offsets []int, err error) {
	if !sort.SliceIsSorted(st.Keys, func(i, j int) bool { return st.Keys[i] < st.Keys[j] }) {
		panic("not sorted!!")
	}
	i := sort.SearchFloat64s(st.Keys, key)
	for ; i < len(st.Keys); i++ {
		if st.Keys[i] > key {
			break
		} else if st.Keys[i] == key {

			offsets = append(offsets, st.Offsets[i])
		}
	}

	if len(offsets) > 0 {
		return offsets, nil
	}
	return nil, fmt.Errorf("The following key <%f> is not found in the index", key)
}

/*
SearchTable uses gonum.sort.Search function to find the leftmost element of the SortedTable
*/
func SearchTable(t Table, x float64) (i int) {
	return sort.Search(
		t.Len(),
		func(i int) bool {
			k, _ := t.Get(i)
			return k >= x
		})
}

/*
InterpolationSearch is a binary search but with the given starting middle
It is used for interpolation, when the middle is not Len/2
*/
func InterpolationSearch(key float64, t Table, m int) (i int) {
	l := 0
	r := t.Len() - 1
	for l < r {
		ikey, _ := t.Get(m)
		if ikey < key {
			l = m + 1
		} else {
			r = m
		}
		m = int(math.Floor(float64(l+r) / 2.))
	}
	return l
}
