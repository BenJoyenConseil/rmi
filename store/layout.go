package store

import (
	"encoding/binary"
	"math"
	"os"
)

const (
	KEY_LEN    = 8         // float64
	VALUE_LEN  = 8         //  uint64
	RECORD_LEN = int64(16) // KEY_LEN + VALUE_LEN
)

// Record is a key/value paire
type Record []byte

func (r Record) Key() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(r[:KEY_LEN]))
}

func (r Record) Value() uint64 {
	return binary.LittleEndian.Uint64(r[VALUE_LEN:])
}

func ToRecord(key float64, value uint64) Record {
	r := make([]byte, KEY_LEN+VALUE_LEN)
	binary.LittleEndian.PutUint64(r[:KEY_LEN], math.Float64bits(key))
	binary.LittleEndian.PutUint64(r[KEY_LEN:], value)
	return r
}

func FromRecord(r Record) (key float64, value uint64) {
	key = math.Float64frombits(binary.LittleEndian.Uint64(r[:KEY_LEN]))
	value = binary.LittleEndian.Uint64(r[KEY_LEN:])
	return key, value
}

// Store is where we store key value paires
type Store struct {
	*os.File
}

func (s Store) Get(i int64) Record {
	offset := i * RECORD_LEN
	b := make([]byte, RECORD_LEN)
	s.ReadAt(b, offset)
	return Record(b)
}

func (s Store) Put(r Record) {
	s.Write(r)
}
