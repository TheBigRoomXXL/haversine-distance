package main

import (
	"fmt"
	"math"
	"testing"
)

func withinTolerance(a, b, e float64) bool {
	d := math.Abs(a - b)
	return d < e
}

// [x0, y0, x1, y1, expectedresult]
var HaversineDistanceCase = [][5]float64{
	{45.544, 50.25, 78.24, 84.21, 3892.318880},
	{98.76, 70.41, -103.43, 12.254, 10670.633119},
	{0, 0, 90, 90, 10010.370831},
	{90, 90, 0, 0, 10010.370831},
	{0, 0, 180, 0, 20020.741663},
	{0, 0, 0, 90, 10010.370831},
	{0, 0, -180, 0, 20020.741663},
	{0, 0, 0, -90, 10010.370831},
}

func TestV0HaversineDistances(t *testing.T) {
	for i, test := range HaversineDistanceCase {
		t.Run(fmt.Sprint(i), func(t *testing.T) {

			result := v0HaversineDistance(test[0], test[1], test[2], test[3], EARTH_RADIUS)
			if !withinTolerance(result, test[4], 0.000001) {
				t.Errorf("want %f, got %f", test[4], result)
			}
		})
	}
}

func TestV1HaversineDistances(t *testing.T) {
	for i, test := range HaversineDistanceCase {
		t.Run(fmt.Sprint(i), func(t *testing.T) {

			result := v1HaversineDistance(test[0], test[1], test[2], test[3], EARTH_RADIUS)
			if !withinTolerance(result, test[4], 0.000001) {
				t.Errorf("want %f, got %f", test[4], result)
			}
		})
	}
}

func benchmarkV0(filename string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		v0(filename)
	}
}

func benchmarkV1(filename string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		v1(filename)
	}
}

func BenchmarkV0_10(b *testing.B)         { benchmarkV0("data/10.json", b) }
func BenchmarkV0_100(b *testing.B)        { benchmarkV0("data/100.json", b) }
func BenchmarkV0_10_000(b *testing.B)     { benchmarkV0("data/10000.json", b) }
func BenchmarkV0_10_000_000(b *testing.B) { benchmarkV0("data/1000000.json", b) }

func BenchmarkV1_10(b *testing.B)         { benchmarkV1("data/10.json", b) }
func BenchmarkV1_100(b *testing.B)        { benchmarkV1("data/100.json", b) }
func BenchmarkV1_10_000(b *testing.B)     { benchmarkV1("data/10000.json", b) }
func BenchmarkV1_10_000_000(b *testing.B) { benchmarkV1("data/1000000.json", b) }
