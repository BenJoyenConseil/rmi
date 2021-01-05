package store

import (
	"encoding/binary"
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
	r := Record([]byte{66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0}) // record(190.223, 4)

	// when
	k, v := FromRecord(r)

	// then
	assert.Equal(t, 190.223, k)
	assert.Equal(t, 4, int(v))
}

func TestKey(t *testing.T) {
	// given
	r := []byte{66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0} // record(190.223, 4)

	// when
	k := Record(r).Key()

	// then
	assert.Equal(t, 190.223, k)
}

func TestValue(t *testing.T) {
	// given
	r := []byte{66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0} // record(190.223, 4)

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
		5, 0, 0, 0, 0, 0, 0, 0, // RecordCount = 5
		66, 96, 229, 208, 34, 199, 103, 64, 1, 0, 0, 0, 0, 0, 0, 0, // Record 0 = record(190.223, 0)
		66, 96, 229, 208, 34, 199, 103, 64, 2, 0, 0, 0, 0, 0, 0, 0, // Record 1 = record(190.223, 1)
		66, 96, 229, 208, 34, 199, 103, 64, 3, 0, 0, 0, 0, 0, 0, 0, // Record 2 = record(190.223, 2)
		66, 96, 229, 208, 34, 199, 103, 64, 4, 0, 0, 0, 0, 0, 0, 0, // Record 3 = record(190.223, 3)
		66, 96, 229, 208, 34, 199, 103, 64, 5, 0, 0, 0, 0, 0, 0, 0, // Record 4 = record(190.223, 4)
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
	r := Record([]byte{1, 0, 0, 8, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0})  // record(1.0, 1)
	r2 := Record([]byte{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0}) // record(2.0, 2)

	// when
	store.Put(r)
	store.Put(r2)

	// then
	result := make([]byte, 16)
	f.ReadAt(result, 8)
	assert.ElementsMatch(t, result, []byte{1, 0, 0, 8, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0})
	f.ReadAt(result, 24)
	assert.ElementsMatch(t, result, []byte{2, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0})

	count := make([]byte, 8)
	f.ReadAt(count, 0)
	assert.ElementsMatch(t, count, []byte{2, 0, 0, 0, 0, 0, 0, 0})
}

func TestRecordCount(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	data := []byte{200, 0, 0, 0, 0, 0, 0, 0}
	f.WriteAt(data, 0)
	store := Store{f}

	// when
	c := store.RecordCount()

	// then
	assert.Equal(t, int64(200), c)
}

func TestRecordCount_NotExits(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	store := Store{f}

	// when
	c := store.RecordCount()

	// then
	assert.Equal(t, int64(0), c)
}

func TestSetRecordCount(t *testing.T) {
	// given
	tmpDir := t.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	store := Store{f}

	// when
	store.setRecordCount(15)

	// then
	f.Seek(0, 0)
	count := make([]byte, 8)
	f.Read(count)
	assert.Equal(t, uint64(15), binary.LittleEndian.Uint64(count))
}

func

func ExampleStore() {

	tmpDir := os.TempDir()
	f, _ := ioutil.TempFile(tmpDir, "*")
	store := Store{f}

	store.Put(ToRecord(1.99, 0))
	store.Put(ToRecord(2.08, 1))
	store.Put(ToRecord(2.08, 2))
	store.Put(ToRecord(2.33, 3))

	fmt.Println("count:", store.RecordCount())
	fmt.Println(FromRecord(store.Get(3)))
	fmt.Println(FromRecord(store.Get(2)))
	fmt.Println(FromRecord(store.Get(1)))
	fmt.Println(FromRecord(store.Get(0)))

	// Output:
	// count: 4
	// 2.33 3
	// 2.08 2
	// 2.08 1
	// 1.99 0
}
