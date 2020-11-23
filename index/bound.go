package index

import "math"

func residual(guess, y int) (residual int) {
	return y - guess
}

/*
scale return the CDF value x datasetLen -1 to get back the position in a sortedTable
*/
func scale(cdfVal float64, datasetLen int) int {
	return int(math.Round(cdfVal*float64(datasetLen) - 1))
}
