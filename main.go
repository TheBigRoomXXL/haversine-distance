package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const NB_CLUSTER = 16
const EARTH_RADIUS = 6372.8

type Pair struct {
	X0 float64 `json:"x0"`
	Y0 float64 `json:"y0"`
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
}

func square(f float64) float64 {
	return f * f
}

func degreeToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func GenerateDataset(n uint64) {
	data := make([]Pair, n)
	result := float64(0)

	if n%10 != 0 {
		log.Fatal("N must be a multiple of 10")
	}

	chunk := n / 10
	for i := uint64(0); i < n; i += chunk {
		x0 := rand.Float64()*350 - 170
		y0 := rand.Float64()*170 - 80
		x1 := rand.Float64()*350 - 170
		y1 := rand.Float64()*170 - 80

		for j := uint64(0); j < chunk; j++ {
			xoffset0 := rand.Float64() * 10
			yoffset0 := rand.Float64() * 10
			xoffset1 := rand.Float64() * 10
			yoffset1 := rand.Float64() * 10

			pair := Pair{
				x0 + xoffset0,
				y0 + yoffset0,
				x1 + xoffset1,
				y1 + yoffset1,
			}
			result += v0HaversineDistance(pair, EARTH_RADIUS)
			data[i+j] = pair
		}
	}

	// Save data to JSON
	file, err := os.Create("data/" + strconv.Itoa(int(n)) + ".json")
	if err != nil {
		log.Fatal("could not create a file to save json data:", err)
	}
	defer file.Close()

	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("could not serialize generated data to json:", err)
	}

	_, err = file.Write(output)
	if err != nil {
		log.Fatal("could not save data to file:", err)
	}
	file.Sync()

	// Save result separatly
	fileResult, err := os.Create("data/" + strconv.Itoa(int(n)))
	if err != nil {
		log.Fatal("could not create a file to save json data:", err)
	}
	defer fileResult.Close()

	result = result / float64(n)
	resultString := fmt.Sprintf("%g", result)

	_, err = fileResult.WriteString(resultString)
	if err != nil {
		log.Fatal("could not save data to file:", err)
	}
	fileResult.Sync()

}

func main() {

	if len(os.Args) < 2 {
		log.Fatal("expected 'generate' or 'process' subcommands")
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
				log.Fatal("bad usage: the number of points to generate is expected")
			}
			GenerateDataset(n)
			return
		}
		GenerateDataset(uint64(1_000_000))
	case "process":
		if len(os.Args) != 4 {
			log.Fatal("bad usage: the implementation and the filepath to the data are expected")
		}

		processor, ok := processors[os.Args[2]]
		if !ok {
			log.Fatal("bad usage: unknown implementation")
		}
		start := time.Now()
		result, n := processor(os.Args[3])
		perf := time.Since(start)
		fmt.Println("Result:", result, "km")

	default:
		log.Fatal("bad usage: 'generate' or 'process' subcommands are expected")
	}
}

var processors = map[string]func(string) (float64, int){
	"v0": v0,
	"v1": v1,
}

var distances = map[string]func(Pair, float64) float64{
	"v0": v0HaversineDistance,
	"v1": v1HaversineDistance,
}
