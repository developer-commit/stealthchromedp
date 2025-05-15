package util

import (
	"math/rand"
)

func NoiseGenerator(value float64, noiseRange float64) float64 {
	return value + (rand.Float64()*2-1)*noiseRange
}
