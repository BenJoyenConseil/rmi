RMI
===
![Go](https://github.com/BenJoyenConseil/learned_index/workflows/Go/badge.svg)

A goland implementation of a RMI ([Recursive Model Indexes](https://github.com/learnedsystems/RMI)), 
a Learned Index structure based on the [research work by Kraska & al.](https://arxiv.org/abs/1712.01208) 

![Fig 1 from the Case for Learned Index Structures](http://people.csail.mit.edu/ryanmarcus/rmi.png)

## usage

Create an index and make lookups

	// load the age column and parse values into float64 values
	ageColumn := extractAgeColumn("data/people.csv")

	// create an index over the age column
	index := index.New(ageColumn)

	// search an age and get back its line position inside the file people.csv
	search, _ := strconv.ParseFloat(os.Args[1], 64)
	line, _ := index.Lookup(search)

the `main.go` file contains an example of a learned index over`data/people.csv` age column. 

It outputs : 

    go run main.go 45
    2020/11/11 15:52:43 Values to index: [90 23 3 45 1 1.5]
    2020/11/11 15:52:43 The value 45 is located line n°5 inside people.csv 

This is the plot showing the approximation (the linear regression), the cumulative distribution function for each value,
and the current age's value (the Keys of the index) : 

![Fig 2 the LearnedIndex over people.csv](assets/plot.svg)

## features

- [x] A simple linear regression model learning the CDF of a float64 array
- [x] A learned index structure fitted on keys of a collection
- [x] Finding rows id on a CSV file
- [ ] Use max + min error bounding elements to search quickly
- [ ] Get all occurences of a key in a dataset
- [ ] Learn on integer
- [ ] A two layer recursive index
- [ ] Index is persistent and durable (on hard drive)
- [ ] A sort algorythm using learned structure
- [ ] Learning on string type ?

## related works

* Tim Kraska, Alex Beutel, Ed H. Chi, Jeffrey Dean, and Neoklis Polyzotis. 2018. The Case for Learned Index Structures. In Proceedings of the 2018 International Conference on Management of Data (SIGMOD '18). Association for Computing Machinery, New York, NY, USA, 489–504. DOI:https://doi.org/10.1145/3183713.3196909

* [Ryan Marcus RMI's reference implementation](https://github.com/learnedsystems/RMI)
