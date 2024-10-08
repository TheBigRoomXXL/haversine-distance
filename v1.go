package main

import (
	"bufio"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

var FLOAT_CHAR = [12]string{
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

func v1IsFloatChar(a string) bool {
	for _, b := range FLOAT_CHAR {
		if b == a {
			return true
		}
	}
	return false
}

func v1JsonToData(reader io.Reader) []Pair {
	data := []Pair{}
	stack := [4]float64{}
	i := uint8(0)

	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanBytes)

	for scanner.Scan() {
		text := scanner.Text()
		if !v1IsFloatChar(text) {
			if text == "x" || text == "y" {
				scanner.Scan()
			}
			continue
		}

		for scanner.Scan() {
			if !v1IsFloatChar(scanner.Text()) {
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
				X0: stack[0],
				Y0: stack[1],
				X1: stack[2],
				Y1: stack[3],
			})
			i = 0
		}
	}

	return data
}

// Same as V0
func v1HaversineDistance(pair Pair, radius float64) float64 {
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

func v1(filepath string) (float64, int) {
	timings := Timings{name: "V1"}
	defer timings.Print()

	timings.Step("open")
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("failed to open the data file")
	}
	defer file.Close()

	timings.Step("parse")
	data := v1JsonToData(file)

	timings.Step("compute")
	sum := 0.0
	for i := 0; i < len(data); i++ {
		dist := v0HaversineDistance(data[i], EARTH_RADIUS)
		sum += dist

	}
	result := sum / float64(len(data))

	timings.N = len(data)
	timings.Step("end")
	return result, len(data)
}
