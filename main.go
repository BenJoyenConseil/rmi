package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"rmi/index"
	"strconv"
)

const FIRST_LINE_OF_DATA int = 2

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Usage: main.go <search_age>")
	}

	// load the age column and parse values into float64 values
	ageColumn := extractAgeColumn("data/people.csv")
	log.Println("Values to index:", ageColumn)

	// create an index over the age column
	index := index.New(ageColumn)
	search, _ := strconv.ParseFloat(os.Args[1], 64)

	// search an age and get back its line position inside the file people.csv
	line, _ := index.Lookup(search)
	log.Printf("The value %s is located line n°%d inside people.csv \n", os.Args[1], line+FIRST_LINE_OF_DATA)
}

func extractAgeColumn(file string) []float64 {
	csvfile, _ := os.Open(file)
	r := csv.NewReader(csvfile)

	var ageColumn []float64
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
				if c == "age" {
					ageCid = i
				}
			}
			headerLine = false
			continue
		}
		age, _ := strconv.ParseFloat(record[ageCid], 64)
		ageColumn = append(ageColumn, age)
	}
	return ageColumn
}
