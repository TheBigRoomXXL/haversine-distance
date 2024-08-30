package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
)

func v0HaversineDistance(pair Pair, radius float64) float64 {
	lat0 := pair.Y0
	lat1 := pair.Y1
	lon0 := pair.X0
	lon1 := pair.X1

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

func v0(filepath string) (float64, int) {
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

	data := []Pair{}
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		fmt.Println("failed to parse JSON", err)
		os.Exit(1)
	}

	sum := 0.0
	for i := 0; i < len(data); i++ {
		dist := v0HaversineDistance(data[i], EARTH_RADIUS)
		sum += dist

	}
	result := sum / float64(len(data))
	return result, len(data)
}
