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

func TestBasicHaversineDistances(t *testing.T) {
	// [x0, y0, x1, y1, expectedresult]
	var tests = [][5]float64{
		{0, 0, 90, 90, 10010.370831},
		{90, 90, 0, 0, 10010.370831},
		{0, 0, 180, 0, 20020.741663},
		{0, 0, 0, 90, 10010.370831},
		{0, 0, -180, 0, 20020.741663},
		{0, 0, 0, -90, 10010.370831},
		{45.544, 50.25, 78.24, 84.21, 3892.318880},
		{98.76, 70.41, -103.43, 12.254, 10670.633119},
	}

	for i, test := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {

			result := HaversineDistance(test[0], test[1], test[2], test[3], EARTH_RADIUS)
			if !withinTolerance(result, test[4], 0.001) {
				t.Errorf("want %f, got %f", test[4], result)
			}
		})
	}
}
