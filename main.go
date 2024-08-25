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

func HaversineDistance(
	x0 float64, y0 float64, x1 float64, y1 float64, radius float64,
) float64 {
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

	if n%10 != 0 {
		panic("N must be a multiple of 10")
	}

	chunk := n / 10
	for i := 0; i < 10; i += 1 {
		x0 := rand.Float64()*350 - 170
		y0 := rand.Float64()*170 - 80
		x1 := rand.Float64()*350 - 170
		y1 := rand.Float64()*170 - 80

		for j := 0; j < chunk; j++ {
			xoffset0 := rand.Float64() * 10
			yoffset0 := rand.Float64() * 10
			xoffset1 := rand.Float64() * 10
			yoffset1 := rand.Float64() * 10

			data[i+j] = [4]float64{
				x0 + xoffset0,
				y0 + yoffset0,
				x1 + xoffset1,
				y1 + yoffset1,
			}
		}
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
	n := 1_000
	defer timer("main", n)()
	data := generateData(n)
	sum := 0.0
	for i := 0; i < len(data); i++ {
		fmt.Println(i, data[i])
		dist := HaversineDistance(
			data[i][0], data[i][1], data[i][2], data[i][3], EARTH_RADIUS,
		)
		sum += dist

	}
	result := sum / float64(len(data))
	fmt.Println("Result:", result)
}
