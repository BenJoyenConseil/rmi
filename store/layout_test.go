package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToRecord(t *testing.T) {
	// when
	b := ToRecord(190.223, 4)

	// then
	assert.ElementsMatch(t, b, []byte{66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0})
}

func TestFromRecord(t *testing.T) {
	// given
	r := Record([]byte{66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0})

	// when
	k, v := FromRecord(r)

	// then
	assert.Equal(t, 190.223, k)
	assert.Equal(t, 4, int(v))
}

func TestKey(t *testing.T) {
	// given
	r := []byte{66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0}

	// when
	k := Record(r).Key()

	// then
	assert.Equal(t, 190.223, k)
}

func TestValue(t *testing.T) {
	// given
	r := []byte{66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0}

	// when
	v := Record(r).Value()

	// then
	assert.Equal(t, 4, int(v))
}

func TestGet(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	data := []byte{
		66, 96, 229, 208, 34, 199, 103, 64, 1, 0, 0, 0, 0, 0, 0, 0,
		66, 96, 229, 208, 34, 199, 103, 64, 2, 0, 0, 0, 0, 0, 0, 0,
		66, 96, 229, 208, 34, 199, 103, 64, 3, 0, 0, 0, 0, 0, 0, 0,
		66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0,
		66, 96, 229, 208, 34, 199, 103, 64, 5, 0, 0, 0, 0, 0, 0, 0,
	}
	f.WriteAt(data, 0)

	store := Store{f}

	// when
	r := store.Get(4)

	// then
	assert.Equal(t, r.Key(), 190.223)
	assert.Equal(t, int(r.Value()), 5)
}

func TestPut(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	store := Store{f}
	r := Record([]byte{66, 96, 229, 208, 34, 199, 103, 64, 5, 0, 0, 0, 0, 0, 0, 0})

	// when
	store.Put(r)

	// then
	result := make([]byte, 16)
	f.ReadAt(result, 0)
	assert.ElementsMatch(t, result, []byte{66, 96, 229, 208, 34, 199, 103, 64, 5, 0, 0, 0, 0, 0, 0, 0})
}

func ExampleStore() {

	tmpDir := os.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	store := Store{f}

	store.Put(ToRecord(1.99, 0))
	store.Put(ToRecord(2.08, 1))
	store.Put(ToRecord(2.08, 2))
	store.Put(ToRecord(2.33, 3))

	fmt.Println(FromRecord(store.Get(3)))
	fmt.Println(FromRecord(store.Get(2)))
	fmt.Println(FromRecord(store.Get(1)))
	fmt.Println(FromRecord(store.Get(0)))

	// Output:
	// 2.33 3
	// 2.08 2
	// 2.08 1
	// 1.99 0
}
