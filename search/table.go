package search

/*
Table is asbtracting behaviors of a sorted collection that has get and insert functions
Implementations of this are InMemory or OnDisk (store)
*/
type Table interface {
	Get(i int) (key float64, rowid int)
	//GetBatch(start, end int) (keys []float64, rowids []int)

	Len() int
}

/*
Searchable gives to structures that implement it the capibility
to return a range of rows matching the key
*/
type Searchable interface {
	Lookup(key float64) []int
}
