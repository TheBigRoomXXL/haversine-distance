package main

import (
	"fmt"
	"math"
	"testing"
)

type testCase struct {
	pair   Pair
	result float64
}

func withinTolerance(a, b, e float64) bool {
	d := math.Abs(a - b)
	return d < e
}

// [x0, y0, x1, y1, expectedresult]
var HaversineDistanceTests = []testCase{
	{Pair{45.544, 50.25, 78.24, 84.21}, 3892.318880},
	{Pair{98.76, 70.41, -103.43, 12.254}, 10670.633119},
	{Pair{0, 0, 90, 90}, 10010.370831},
	{Pair{90, 90, 0, 0}, 10010.370831},
	{Pair{0, 0, 180, 0}, 20020.741663},
	{Pair{0, 0, 0, 90}, 10010.370831},
	{Pair{0, 0, -180, 0}, 20020.741663},
	{Pair{0, 0, 0, -90}, 10010.370831},
}

func TestHaversineDistances(t *testing.T) {
	for i, test := range HaversineDistanceTests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result := v0HaversineDistance(test.pair, EARTH_RADIUS)
			if !withinTolerance(result, test.result, 0.000001) {
				t.Errorf("want %f, got %f", test.result, result)
			}
		})
	}
}

func TestV1HaversineDistances(t *testing.T) {
	for i, test := range HaversineDistanceTests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {

			result := v1HaversineDistance(test.pair, EARTH_RADIUS)
			if !withinTolerance(result, test.result, 0.000001) {
				t.Errorf("want %f, got %f", test.result, result)
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

func BenchmarkV0_200(b *testing.B)       { benchmarkV0("data/200.json", b) }
func BenchmarkV0_40_000(b *testing.B)    { benchmarkV0("data/40000.json", b) }
func BenchmarkV0_1_000_000(b *testing.B) { benchmarkV0("data/1000000.json", b) }

func BenchmarkV1_200(b *testing.B)       { benchmarkV1("data/200.json", b) }
func BenchmarkV1_40_000(b *testing.B)    { benchmarkV1("data/40000.json", b) }
func BenchmarkV1_1_000_000(b *testing.B) { benchmarkV1("data/1000000.json", b) }
