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
	  search, _ := strconv.ParseFloat(os.Args[1], 64)

    // search an age and get back its line position inside the file people.csv
    line, _ := index.Lookup(search)

the `main.go` file contains an example of a learned index over`data/people.csv` age column. 

It outputs : 

    go run main.go 45
    2020/11/11 15:52:43 Values to index: [90 23 3 45 1 1.5]
    2020/11/11 15:52:43 The value 45 is located line n°5 inside people.csv 

## features

- [x] A learned index structure using linear regression fitted on keys of type float
- [ ] learning on integer
- [ ] a two layer recursive index
- [ ] a sort algorythm using learned structure
- [ ] recursive indexes with topdown approarch
- [ ] learning on string type

## related works

* Tim Kraska, Alex Beutel, Ed H. Chi, Jeffrey Dean, and Neoklis Polyzotis. 2018. The Case for Learned Index Structures. In Proceedings of the 2018 International Conference on Management of Data (SIGMOD '18). Association for Computing Machinery, New York, NY, USA, 489–504. DOI:https://doi.org/10.1145/3183713.3196909

* [Ryan Marcus RMI's reference implementation](https://github.com/learnedsystems/RMI)
