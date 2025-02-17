package util

import (
	"math"
	"math/rand"
)

func DegressToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func RangeRange(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}
