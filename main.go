package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"rmi/index"
	"strconv"
	"strings"
)

const FIRST_LINE_OF_DATA int = 2

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Usage: main.go <search_age>")
	}
	file := "data/titanic.csv"
	// load the age column and parse values into float64 values
	ageColumn := extractColumn(file, "age")

	// create an index over the age column
	idx := index.New(ageColumn)
	search, _ := strconv.ParseFloat(os.Args[1], 64)

	// search an age and get back its line position inside the file people.csv
	result, err := idx.Lookup(search)
	if err != nil {
		log.Fatalf("There is no entry found for %s inside %s \n", os.Args[1], file)
	}
	lines := []int{}
	for _, l := range result {
		lines = append(lines, l+FIRST_LINE_OF_DATA)
	}
	log.Printf("People who are %s years old are located at %d inside %s \n", os.Args[1], lines, file)
	png, _ := filepath.Abs("assets/plot.svg")
	index.Genplot(idx, ageColumn, png)
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
