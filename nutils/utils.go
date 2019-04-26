package nutils

import "math"

func IntMax(a ...int) int {
	m := math.MinInt32
	for _, v := range a {
		if m < v {
			m = v
		}
	}
	return m
}

