package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
)

func baseline(filepath string) (float64, int) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("failed to open the data file")
		os.Exit(1)
	}

	jsonBytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("failed to read the data file", err)
		os.Exit(1)
	}

	data := [][4]float64{}
	json.Unmarshal(jsonBytes, &data)

	sum := 0.0
	for i := 0; i < len(data); i++ {
		dist := baselineHaversineDistance(
			data[i][0], data[i][1], data[i][2], data[i][3], EARTH_RADIUS,
		)
		sum += dist

	}
	result := sum / float64(len(data))
	return result, len(data)
}

func baselineHaversineDistance(
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

var implementations = map[string]func(string) (float64, int){
	"baseline": baseline,
}
