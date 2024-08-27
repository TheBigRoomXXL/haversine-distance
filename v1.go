package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
)

type Pair struct {
	x0 float64
	y0 float64
	x1 float64
	y1 float64
}

var FLORAT_CHAR = [12]string{
	"-",
	".",
	"0",
	"1",
	"2",
	"3",
	"4",
	"5",
	"6",
	"7",
	"8",
	"9",
}

func isFloatChar(a string) bool {
	for _, b := range FLORAT_CHAR {
		if b == a {
			return true
		}
	}
	return false
}

func jsonToData(reader io.Reader) []Pair {
	data := []Pair{}
	stack := [4]float64{}
	i := uint8(0)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		text := scanner.Text()
		if !isFloatChar(text) {
			continue
		}

		for scanner.Scan() {
			if !isFloatChar(scanner.Text()) {
				break
			}
			text += scanner.Text()
		}
		value, err := strconv.ParseFloat(text, 64)
		if err != nil {
			panic(err)
		}

		stack[i] = value
		i++

		if i == 4 {
			data = append(data, Pair{
				x0: stack[0],
				y0: stack[1],
				x1: stack[2],
				y1: stack[3],
			})
			i = 0
		}
	}

	return data
}

func v1HaversineDistance(
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

func v1(filepath string) (float64, int) {
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("failed to open the data file")
		os.Exit(1)
	}

	data := jsonToData(file)

	sum := 0.0
	for i := 0; i < len(data); i++ {
		dist := v0HaversineDistance(
			data[i].x0, data[i].y0, data[i].x1, data[i].y1, EARTH_RADIUS,
		)
		sum += dist

	}
	result := sum / float64(len(data))
	return result, len(data)
}
