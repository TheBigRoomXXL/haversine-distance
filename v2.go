package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

var FLOAT_BYTE = []byte("-.0123456789")

func v2IsFloatByte(a byte) bool {
	for _, b := range FLOAT_BYTE {
		if b == a {
			return true
		}
	}
	return false
}

func v2JsonToData(reader io.Reader) []Pair {
	data := []Pair{}
	stack := [4]float64{}
	i := uint8(0)

	br := bufio.NewReader(reader)

	for {
		text := []byte{}
		b, err := br.ReadByte()
		if err != nil && errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			panic(err)
		}

		if !v2IsFloatByte(b) {
			if b == byte(120) || b == byte(121) {
				_, err := br.ReadByte()
				if err != nil {
					panic(err)
				}
			}
			continue
		}

		text = append(text, b)

		for {
			b, err := br.ReadByte()
			if err != nil && !errors.Is(err, io.EOF) {
				panic(err)
			}
			if !v2IsFloatByte(b) {
				break
			}
			text = append(text, b)
		}

		value, err := strconv.ParseFloat(string(text), 64)
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
func v2HaversineDistance(pair Pair, radius float64) float64 {
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

func v2(filepath string) (float64, int) {
	timings := Timings{name: "V2"}
	defer timings.Print()

	timings.Step("open")
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("failed to open the data file")
	}
	defer file.Close()

	timings.Step("parse")
	data := v2JsonToData(file)

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
