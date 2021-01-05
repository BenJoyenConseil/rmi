package store

import (
	"encoding/binary"
	"log"
	"math"
	"os"
)

const (
	KEY_LEN    = 8         // float64
	VALUE_LEN  = 8         //  uint64
	RECORD_LEN = int64(16) // KEY_LEN + VALUE_LEN
	HEADER     = 8

	VERSION          = "v0.0.0"
	METADATA_LEN     = 4
	RECORD_COUNT_LEN = 8
	HEADERS          = len(VERSION) + RECORD_COUNT_LEN + METADATA_LEN
)

// Record is a key/value paire
type Record []byte

func (r Record) Key() float64 {
	return math.Float64frombits(binary.LittleEndian.Uint64(r[:KEY_LEN]))
}

func (r Record) Value() uint64 {
	return binary.LittleEndian.Uint64(r[KEY_LEN:])
}

func ToRecord(key float64, value uint64) Record {
	r := make([]byte, RECORD_LEN)
	binary.LittleEndian.PutUint64(r[:KEY_LEN], math.Float64bits(key))
	binary.LittleEndian.PutUint64(r[KEY_LEN:], value)
	return r
}

func FromRecord(r Record) (key float64, value uint64) {
	key = math.Float64frombits(binary.LittleEndian.Uint64(r[:KEY_LEN]))
	value = binary.LittleEndian.Uint64(r[KEY_LEN:])
	return key, value
}

// Store is a os.File where we store key value paires
type Store struct {
	*os.File
}

/*
Get reads the store file at offset i and return a Record byte array
*/
func (s Store) Get(i int64) Record {
	offset := i*RECORD_LEN + HEADER
	b := make([]byte, RECORD_LEN)
	_, err := s.ReadAt(b, offset)
	check(err)
	return Record(b)
}

/*
Put appends a Record to the store file
*/
func (s Store) Put(r Record) {
	count := s.RecordCount()
	offset := count*RECORD_LEN + HEADER
	log.Println(count, RECORD_LEN, HEADER, offset)
	_, err := s.WriteAt(r, offset)
	check(err)
	s.setRecordCount(count + 1)
}

/*
RecordCount reads the first byte as
*/
func (s Store) RecordCount() int64 {
	b := make([]byte, HEADER)
	_, err := s.ReadAt(b, 0)
	if err != nil {
		return int64(0)
	}
	return int64(binary.LittleEndian.Uint64(b))
}

func (s Store) setRecordCount(n int64) {
	b := make([]byte, HEADER)
	binary.LittleEndian.PutUint64(b, uint64(n))

	_, err := s.WriteAt(b, 0)
	check(err)
}

func (s Store) MetadataLen() int32 {
	panic("not implemented")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
