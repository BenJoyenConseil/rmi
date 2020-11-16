package index

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"rmi/linear"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSortedTable(t *testing.T) {
	// given
	unsortedkeys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	st := newSortedTable(unsortedkeys)

	// then
	assert.Equal(t, []float64{2.5, 2.98, 3, 3, 3.14, 5, 10}, st.keys)
	assert.Equal(t, []int{5, 6, 1, 2, 3, 0, 4}, st.offsets)
}
func TestNew(t *testing.T) {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	// when
	idx := New(keys)

	// then
	assert.Equal(t, 7, idx.Len)
	assert.Equal(t, .23119036646681634, idx.M.Intercept)
	assert.Equal(t, .08523040437506509, idx.M.Slope)
	assert.Equal(t, 2, idx.MaxErrBound)
	assert.Equal(t, -2, idx.MinErrBound)
	// check the order
	assert.Equal(t, []float64{2.5, 2.98, 3, 3, 3.14, 5, 10}, idx.ST.keys)
	assert.Equal(t, []int{5, 6, 1, 2, 3, 0, 4}, idx.ST.offsets)
}

func TestGuessIndex(t *testing.T) {
	idx := &LearnedIndex{
		M:           &linear.RegressionModel{Intercept: .23119036646681634, Slope: .08523040437506509},
		Len:         7,
		MaxErrBound: 2,
		MinErrBound: -3,
	}

	// when
	guess, lower, upper := idx.GuessIndex(2.95)
	assert.Equal(t, 2, guess)
	assert.Equal(t, 0, lower)
	assert.Equal(t, 4, upper)

	// when guess < 0 ==> -1
	guess, lower, upper = idx.GuessIndex(-2.95)
	assert.Equal(t, 0, guess)
	assert.Equal(t, 0, lower)
	assert.Equal(t, 2, upper)

	// when guess > len-1 ==> 7
	guess, lower, upper = idx.GuessIndex(10.)
	assert.Equal(t, 6, guess)
	assert.Equal(t, 3, lower)
	assert.Equal(t, 6, upper)
}

func TestLookup(t *testing.T) {
	// given
	idx := &LearnedIndex{
		M:           &linear.RegressionModel{Intercept: .23119036646681634, Slope: .08523040437506509},
		Len:         7,
		MaxErrBound: 2,
		MinErrBound: -2,
		//	keys : {5, 3, 3, 3.14, 10, 2.5, 2.98}
		//  sort : {2.5, 2.98, 3, 3, 3.14, 5, 10}
		//  posi : {0,   1,    2, 3, 4,    5, 6}
		ST: &sortedTable{
			keys:    []float64{2.5, 2.98, 3, 3, 3.14, 5, 10},
			offsets: []int{5, 6, 1, 2, 3, 0, 4},
		},
	}

	// when
	offsets, err := idx.Lookup(2.5)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{5}, offsets)
	// when
	offsets, err = idx.Lookup(2.98)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{6}, offsets)
	// when
	offsets, err = idx.Lookup(3.)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{1, 2}, offsets)
	// when
	offsets, err = idx.Lookup(3.14)
	assert.Nil(t, err)
	assert.ElementsMatch(t, []int{3}, offsets)
	// when
	offsets, err = idx.Lookup(5.)
	assert.Nil(t, err)
	assert.Equal(t, []int{0}, offsets)
	// when
	offsets, err = idx.Lookup(10.)
	assert.Nil(t, err)
	assert.Equal(t, []int{4}, offsets)

	// when not in the index
	offsets, err = idx.Lookup(199.)
	assert.NotNil(t, err)
	assert.Nil(t, offsets)
}

func ExampleLearnedIndex() {
	keys := []float64{5, 3, 3, 3.14, 10, 2.5, 2.98}

	index := New(keys)

	for _, k := range keys {
		offset, err := index.Lookup(k)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf("The key %f is located %d\n", k, offset)
	}

	// Output:
	// The key 2.500000 is located [5]
	// The key 2.980000 is located [6]
	// The key 3.000000 is located [1 2]
	// The key 3.000000 is located [1 2]
	// The key 3.140000 is located [3]
	// The key 5.000000 is located [0]
	// The key 10.000000 is located [4]
}

var min, max = 0., 100.
var random = func() float64 { return math.Round(min + rand.Float64()*(max-min)) }

func BenchmarkLearnedIndex(b *testing.B) {
	file := "../data/titanic.csv"
	// load the age column and parse values into float64 values
	ageColumn := extractColumn(file, "age")

	// create an index over the age column
	idx := New(ageColumn)
	keyFound := map[float64][]int{}
	keyNotFound := map[float64][]error{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := random()
		offsets, err := idx.Lookup(k)
		if err != nil {
			keyNotFound[k] = append(keyNotFound[k], err)
		} else {
			keyFound[k] = offsets
		}
	}
	log.Println("keys found:", len(keyFound), "| keys not found:", len(keyNotFound))
	//log.Println(keyNotFound)
}

func BenchmarkBinarySearch(b *testing.B) {
	file := "../data/titanic.csv"
	// load the age column and parse values into float64 values
	ageColumn := extractColumn(file, "age")
	sort.Float64s(ageColumn)
	keyFound := map[float64][]int{}
	keyNotFound := map[float64][]error{}
	b.ResetTimer()
	o := 0
	for i := 0; i < b.N; i++ {
		k := random()
		o = sort.SearchFloat64s(ageColumn, k)

		if o >= len(ageColumn) {
			keyNotFound[k] = append(keyNotFound[k], fmt.Errorf("%f not found", k))
		} else {
			offsets := []int{}
			for o < len(ageColumn) {
				if ageColumn[o] != k {
					break
				}
				offsets = append(offsets, o)
				o++
			}
			keyFound[k] = offsets
		}
	}
	log.Println("keys found:", len(keyFound), "| keys not found:", len(keyNotFound))
	//log.Println(keyNotFound)
}

func extractColumn(file string, colName string) []float64 {
	csvfile, _ := os.Open(file)
	r := csv.NewReader(csvfile)

	var valuesColumn []float64
	var ageCid int
	var headerLine bool = true
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if headerLine {
			for i, c := range record {
				if strings.ToLower(c) == colName {
					ageCid = i
				}
			}
			headerLine = false
			continue
		}
		v, _ := strconv.ParseFloat(record[ageCid], 64)
		valuesColumn = append(valuesColumn, v)
	}
	return valuesColumn
}
