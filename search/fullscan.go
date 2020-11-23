package search

import "fmt"

func FullScanLookup(key float64, st *SortedTable) (offsets []int, err error) {

	for i := 0; i < len(st.Keys); i++ {
		if st.Keys[i] == key {
			offsets = append(offsets, st.Offsets[i])
		}
	}

	if len(offsets) > 0 {
		return offsets, nil
	}
	return nil, fmt.Errorf("The following key <%f> is not found in the index", key)
}
