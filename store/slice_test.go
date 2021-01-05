package store

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSlice(t *testing.T) {
	//given
	tmpDir := t.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	store := Store{f}

	// when
	slice := NewSlice(store, 0, 3)

	// then
	assert.Equal(t, int64(0), slice.start)
	assert.Equal(t, int64(3), slice.end)
	assert.Equal(t, store, slice.store)
}

func TestNewSliceFail_WhenStartUpperEnd(t *testing.T) {
	//given
	store := Store{}
	start := int64(0)
	end := int64(-3)

	// when
	// then
	assert.Panics(t, func() { NewSlice(store, start, end) })
}

func TestSliceLen(t *testing.T) {
	// given
	slice := &Slice{
		start: 10,
		end:   190,
	}

	// when
	l := slice.Len()
	// then
	assert.Equal(t, int64(181), l)
}

func TestSliceGet(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	data := []byte{
		5, 0, 0, 0, 0, 0, 0, 0, // RecordCount = 5
		66, 96, 229, 208, 34, 199, 103, 64, 1, 0, 0, 0, 0, 0, 0, 0, // Record 0 = record(190.223, 0)
		66, 96, 229, 208, 34, 199, 103, 64, 2, 0, 0, 0, 0, 0, 0, 0, // Record 1 = record(190.223, 1)
		66, 96, 229, 208, 34, 199, 103, 64, 3, 0, 0, 0, 0, 0, 0, 0, // Record 2 = record(190.223, 2)
		66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0, // Record 3 = record(190.223, 3)
		66, 96, 229, 208, 34, 199, 103, 64, 5, 0, 0, 0, 0, 0, 0, 0, // Record 4 = record(190.223, 4)
	}
	f.WriteAt(data, 0)
	store := Store{f}
	slice := &Slice{
		store: store,
		start: 1,
		end:   3,
	}
	i := int64(2)

	// when
	key, rowid := slice.Get(i)

	// then
	assert.Equal(t, 190.223, key)
	assert.Equal(t, uint64(3), rowid)
}
