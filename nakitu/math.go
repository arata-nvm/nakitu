package nakitu

import (
	"math"
	"math/rand"
)

func Rad(deg float64) float64 {
	return deg * math.Pi / 180
}

func Clamp(x, min, max float64) float64 {
	return math.Max(math.Min(x, max), min)
}

func Random(min, max float64) float64 {
	return min + (max-min)*rand.Float64()
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}
