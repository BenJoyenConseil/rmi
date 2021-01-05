package search

import "sort"

/*
A Sorted Table represents a collection of key:offset pairs that is sorted by key
keeping offsets following their corresponding key
*/
type SortedTable struct {
	Keys    []float64
	Offsets []int
}

// Implement the Sort interface
type byKeys struct{ *SortedTable }

func (st byKeys) Len() int { return len(st.Keys) }
func (st byKeys) Swap(i, j int) {
	st.Keys[i], st.Keys[j] = st.Keys[j], st.Keys[i]
	st.Offsets[i], st.Offsets[j] = st.Offsets[j], st.Offsets[i]
}
func (st byKeys) Less(i, j int) bool { return st.Keys[i] < st.Keys[j] }

/*
Return a Sorted Table structure sorted by key in an ascending order
*/
func NewSortedTable(x []float64) *SortedTable {
	keys, offsets := x, make([]int, len(x))
	for i := range x {
		offsets[i] = i
	}
	st := &SortedTable{Keys: keys, Offsets: offsets}
	sort.Sort(byKeys{st})
	return st
}

func (st *SortedTable) Get(i int) (key float64, rowid int) {
	return st.Keys[i], st.Offsets[i]
}

func (st *SortedTable) Len() int {
	return len(st.Keys)
}
