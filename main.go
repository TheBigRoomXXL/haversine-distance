package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const NB_CLUSTER = 16
const EARTH_RADIUS = 6372.8

func square(f float64) float64 {
	return f * f
}

func degreeToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func HaversineDistance(x0 float64, y0 float64, x1 float64, y1 float64, radius float64) float64 {
	lat0 := y0
	lat1 := y1
	lon0 := x0
	lon1 := x1

	dLat := degreeToRadians(lat1 - lat0)
	dLon := degreeToRadians(lon1 - lon0)
	lat0 = degreeToRadians(lat0)
	lat1 = degreeToRadians(lat1)

	a := square(math.Sin(dLat / 2.0))
	b := math.Cos(lat0) * math.Cos(lat1) * square(math.Sin(dLon/2))
	c := 2.0 * math.Asin(math.Sqrt(a+b))
	result := radius * c
	return result
}

func generateData(n int) [][4]float64 {
	data := make([][4]float64, n)

	for i := 0; i < n; i++ {
		x0 := rand.Float64()*360 - 180
		y0 := rand.Float64()*180 - 90
		x1 := rand.Float64()*360 - 180
		y1 := rand.Float64()*180 - 90
		data[i] = [4]float64{x0, y0, x1, y1}
	}
	return data
}

func timer(name string, n int) func() {
	start := time.Now()
	return func() {
		perf := time.Since(start)
		fmt.Printf("⏱ %s took %v for %.d calculations\n", name, perf, n)
		fmt.Printf("  ↳ %f µs/calc\n", float64(perf.Microseconds())/float64(n))
	}
}

func main() {
	n := 1_000_000
	defer timer("main", n)()
	data := generateData(n)
	sum := 0.0
	for i := 0; i < len(data); i++ {
		sum += HaversineDistance(data[i][0], data[i][1], data[i][2], data[i][3], EARTH_RADIUS)
	}
	result := sum / float64(len(data))
	fmt.Println("Result:", result)
}
