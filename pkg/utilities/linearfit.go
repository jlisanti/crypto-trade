package utilities

import (
	"math"
)

func LinearFit(y []float64, x []float64, m float64, c float64) (float64, float64) {
	if len(y) > 1 {
		tol := false
		maxIter := 100000
		tolerance := 0.01
		L := 0.001
		i := 0
		for !tol {
			dm := partialm(y, x, m, c)
			dc := partialc(y, x, m, c)

			m -= L * dm
			c -= L * dc

			computedErr := calcError(y, x, m, c)

			if (computedErr < tolerance) || i > maxIter {
				break
			}
			i += 1
		}
		return m, c
	} else {
		return 0.0, 0.0
	}
}

func partialm(y []float64, x []float64, m float64, c float64) float64 {

	n := len(y)
	sum := 0.0
	for i, yi := range y {
		sum += x[i] * (yi - (x[i]*m + c))
	}
	return -2.0 / float64(n) * sum
}

func partialc(y []float64, x []float64, m float64, c float64) float64 {

	n := len(y)
	sum := 0.0
	for i, yi := range y {
		sum += yi - (x[i]*m + c)
	}
	return -2.0 / float64(n) * sum
}

func calcError(y []float64, x []float64, m float64, c float64) float64 {
	sum := 0.0
	for i, yi := range y {
		sum += math.Pow(yi-(m*x[i]+c), 2)
	}
	return 1.0 / float64(len(y)) * sum
}
