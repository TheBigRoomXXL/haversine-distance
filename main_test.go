package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
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

// ============== TEST DISTANCE CALCULATION ==============

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

func testHaversine(t *testing.T, version string) {
	distanceFunc := distances[version]
	for i, test := range HaversineDistanceTests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			result := distanceFunc(test.pair, EARTH_RADIUS)
			if !withinTolerance(result, test.result, 0.000001) {
				t.Errorf("want %f, got %f", test.result, result)
			}
		})
	}
}

func TestHaversineV0(t *testing.T) { testHaversine(t, "v0") }
func TestHaversineV1(t *testing.T) { testHaversine(t, "v1") }

// ============== TEST PROCESSOR ==============

func testProcessor(t *testing.T, n int, version string) {
	dataFile := "data/" + strconv.Itoa(int(n)) + ".json"
	resultFile := "data/" + strconv.Itoa(int(n))

	b, err := os.ReadFile(resultFile)
	if err != nil {
		log.Fatal(err)
	}

	expected, err := strconv.ParseFloat(string(b), 64)
	if err != nil {
		log.Fatal(err)
	}

	result, nResult := processors[version](dataFile)
	if n != nResult {
		t.Errorf("want %d, got %d", n, nResult)
	}
	if !withinTolerance(expected, result, 0.000001) {
		t.Errorf("want %f, got %f", expected, result)
	}
}

func TestProcessorV0_200(t *testing.T)  { testProcessor(t, 200, "v0") }
func TestProcessorV0_3000(t *testing.T) { testProcessor(t, 3000, "v0") }
func TestProcessorV1_200(t *testing.T)  { testProcessor(t, 200, "v1") }
func TestProcessorV1_3000(t *testing.T) { testProcessor(t, 3000, "v1") }

// ============== BENCHMARKS ===============

func benchmark(b *testing.B, n int, version string) {
	processor := processors[version]
	dataFile := "data/" + strconv.Itoa(int(n)) + ".json"
	for n := 0; n < b.N; n++ {
		processor(dataFile)
	}
}

func Benchmark_V0_200(b *testing.B)    { benchmark(b, 200, "v0") }
func Benchmark_V0_40000(b *testing.B)  { benchmark(b, 40000, "v0") }
func Benchmark_V0_500000(b *testing.B) { benchmark(b, 500000, "v0") }
func Benchmark_V1_200(b *testing.B)    { benchmark(b, 200, "v1") }
func Benchmark_V1_40000(b *testing.B)  { benchmark(b, 40000, "v1") }
func Benchmark_V1_500000(b *testing.B) { benchmark(b, 500000, "v1") }
