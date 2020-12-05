package cli

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/BenJoyenConseil/rmi/index"
	"github.com/BenJoyenConseil/rmi/store"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	IndexFile            = "index.rmi"
	IndexFileName        = "./index.rmi"
	DefaultPlotImgFormat = "svg"
	DefaultPlotFolder    = "assets"
)

var (
	app = kingpin.New("rmi", "learns from data to index efficiently your CSV")

	create        = app.Command("create", "build an index structure that learn distribution over values of a column")
	fileToIndex   = create.Flag("csv", "The CSV file you want to index").Short('f').Required().String()
	columnToIndex = create.Flag("column", "The column you want to index").Short('c').Required().String()
	createAction  = create.Action(createIndex)

	select_         = app.Command("search", "query the index file found in the targeted directory")
	selectIndexFile = select_.Flag("index", "the index file").Short('i').Default(IndexFileName).ExistingFile()
	searchedValue   = select_.Arg("key", "designates the key used to find the corresponding lines").String()
	selectAction    = select_.Action(selectWhere)

	plot                 = app.Command("plot", "print a graphic representation of the index, its cdf, the approximation used")
	plotDir              = plot.Flag("dir", "the Dir where to write the resulting file").Short('d').Default(DefaultPlotFolder).ExistingDir()
	plotExt              = plot.Flag("type", "The image type : png, svg, jpg").Short('t').Default(DefaultPlotImgFormat).Enum("svg", "png", "jpg")
	plotSmoothBoundaries = plot.Flag("smooth", "define if the boundaries are smoothed or let them raw").Default("true").Bool()
	plotAction           = plot.Action(plotIndex)
)

func Parse(args []string) {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func createIndex(c *kingpin.ParseContext) error {
	// file := "data/people.csv"
	ageColumn := extractColumn(*fileToIndex, *columnToIndex)

	// create an index over the age column
	idx := index.New(ageColumn)
	storeFile, _ := os.OpenFile(IndexFileName, os.O_CREATE|os.O_WRONLY, 0644)
	s := store.Store{storeFile}
	for i := 0; i < idx.Len; i++ {
		s.Put(store.ToRecord(idx.ST.Keys[i], uint64(idx.ST.Offsets[i])))
	}
	return nil
}

func selectWhere(c *kingpin.ParseContext) error {
	/*
		const FIRST_LINE_OF_DATA int = 2
		search, _ := strconv.ParseFloat(os.Args[1], 64)
		log.Println("max error is :", idx.MaxErrBound, "; min error is", idx.MinErrBound)

		// search an age and get back its line position inside the file people.csv
		result, err := idx.Lookup(search)
		if err != nil {
			log.Fatalf("There is no entry found for %s inside %s \n", os.Args[1], file)
		}
		lines := []int{}
		for _, l := range result {
			lines = append(lines, l+FIRST_LINE_OF_DATA)
		}
		log.Printf("We found %d entries in the index \n", len(lines))
		log.Printf("People who are %s years old are located at %d inside %s \n", os.Args[1], lines, file)
	*/
	return nil
}

func plotIndex(c *kingpin.ParseContext) error {
	// ageColumn := extractColumn(file, "age")
	//img, _ := filepath.Abs(*plotPath)
	//index.Genplot(idx, ageColumn, img, *plotAction)
	return nil
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
