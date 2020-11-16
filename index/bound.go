package index

import "math"

func residual(guesses, y []int) (min, max int) {
	for i := range guesses {
		residual := y[i] - guesses[i]
		if residual < min {
			min = residual
		}
		if residual > max {
			max = residual
		}
	}
	return min, max
}

/*
scale return the CDF value x datasetLen -1 to get back the position in a sortedTable
*/
func scale(cdfVal float64, datasetLen int) int {
	return int(math.Round(cdfVal*float64(datasetLen) - 1))
}
