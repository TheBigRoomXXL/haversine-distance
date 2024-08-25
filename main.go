package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
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

func GenerateDataset(n uint64) [][4]float64 {
	data := make([][4]float64, n)

	if n%10 != 0 {
		fmt.Println("N must be a multiple of 10")
		os.Exit(1)
	}

	chunk := n / 10
	for i := uint64(0); i < 10; i += 1 {
		x0 := rand.Float64()*350 - 170
		y0 := rand.Float64()*170 - 80
		x1 := rand.Float64()*350 - 170
		y1 := rand.Float64()*170 - 80

		for j := uint64(0); j < chunk; j++ {
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

	file, err := os.Create("data/" + strconv.Itoa(int(n)) + ".json")
	if err != nil {
		fmt.Println("could not create a file to save json data:", err)
		os.Exit(1)
	}
	defer file.Close()

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("could not serialize generated data to json:", err)
		os.Exit(1)
	}

	_, err = file.Write(output)
	if err != nil {
		fmt.Println("could not save data to file:", err)
		os.Exit(1)
	}
	file.Sync()

	return data
}

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
		dist := HaversineDistance(
			data[i][0], data[i][1], data[i][2], data[i][3], EARTH_RADIUS,
		)
		sum += dist

	}
	result := sum / float64(len(data))
	return result, len(data)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("expected 'generate' or 'process' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--help":
		fmt.Println("TODO usage")

	case "-h":
		fmt.Println("TODO usage")

	case "generate":
		if len(os.Args) > 2 {
			n, err := strconv.ParseUint(os.Args[2], 10, 64)
			if err != nil {
				fmt.Println("bad usage: the number of points to generate is expected")
			}
			GenerateDataset(n)
			return
		}
		GenerateDataset(uint64(1_000_000))
	case "process":
		if len(os.Args) != 4 {
			fmt.Println("bad usage: the implementation and the filepath to the data are expected")
			os.Exit(1)
		}

		processor, ok := implementations[os.Args[2]]
		if !ok {
			fmt.Println("bad usage: unknown implementation")
			os.Exit(1)
		}
		start := time.Now()
		result, n := processor(os.Args[3])
		perf := time.Since(start)
		fmt.Println("Result:", result, "km")
		fmt.Printf("⏱ %s took %v for %.d calculations\n", os.Args[2], perf, n)
		fmt.Printf("  ↳ %f µs/calc\n", float64(perf.Microseconds())/float64(n))

	default:
		fmt.Println("bad usage: 'generate' or 'process' subcommands are expected")
		os.Exit(1)
	}

}

var implementations = map[string]func(string) (float64, int){
	"baseline": baseline,
}
