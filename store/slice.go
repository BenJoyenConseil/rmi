package store

import "fmt"

type Slice struct {
	store Store
	start int64
	end   int64
}

func NewSlice(s Store, start, end int64) *Slice {
	if end < start {
		panic(fmt.Sprintf("The end :%d cannot be lower than the start :%d", end, start))
	}
	return &Slice{
		store: s,
		start: start,
		end:   end,
	}
}

func (slice *Slice) Len() int64 {
	return slice.end - slice.start + 1
}

func (slice *Slice) Get(i int64) (float64, uint64) {
	return FromRecord(slice.store.Get(i))
}
