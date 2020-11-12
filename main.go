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

	// load the age column and parse values into float64 values
	ageColumn := extractColumn("data/titanic.csv", "age")
	log.Println("Values to index:", ageColumn)

	// create an index over the age column
	idx := index.New(ageColumn)
	search, _ := strconv.ParseFloat(os.Args[1], 64)

	// search an age and get back its line position inside the file people.csv
	line, err := idx.Lookup(search)
	if err != nil {
		log.Fatalf("The value %s is not found inside people.csv \n", os.Args[1])
	}
	log.Printf("The value %s is located line n°%d inside people.csv \n", os.Args[1], line+FIRST_LINE_OF_DATA)
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
