package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/BenJoyenConseil/rmi/index"
	"github.com/BenJoyenConseil/rmi/search"
	"github.com/stretchr/testify/assert"
)

func TestIsoFunctional(t *testing.T) {

	// given the titanic.csv dataset
	ageCol := extractColumn("./data/titanic.csv", "age")
	li := index.New(ageCol)
	log.Println(li.MaxErrBound, li.MinErrBound)
	//ci := NewCubic(ageCol)

	// when Lookup using bsearch, learnedindex, etc..
	for i := 0.; i <= 100; i++ {

		resultFS, errFS := search.FullScanLookup(i, li.ST)
		resultLI, errLI := li.Lookup(i)
		//resultCI, errCI := ci.Lookup(i)

		// then forearch key result should be the same
		assert.ElementsMatch(t, resultFS, resultLI, i)
		//assert.ElementsMatch(t, resultBS, resultCI)
		assert.Equal(t, errFS, errLI, i)
		//assert.Equal(t, errBS, errCI)
	}
}

var min, max = 0., 100.
var random = func() float64 { return math.Round(min + rand.Float64()*(max-min)) }

func BenchmarkLearnedIndex(b *testing.B) {
	file := "./data/titanic.csv"
	// load the age column and parse values into float64 values
	ageColumn := extractColumn(file, "age")

	// create an index over the age column
	idx := index.New(ageColumn)
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
	file := "./data/titanic.csv"
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
