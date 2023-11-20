package preprocess

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Loads map[int]*LoadSingle

type LoadSingle struct {
	Pickup  *Coordinates
	Dropoff *Coordinates
}

type Coordinates struct {
	X float64
	Y float64
}

func GetData() Loads {
	// get args (filepath)
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Usage: go run main.go {path_to_problem} or main {path_to_problem}")
	}

	// access the file
	filePath := args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error reading the file at: ", filePath)
	}

	// create loadsData that is to be returned
	loadsData := make(Loads)

	// read the file line by line
	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		newLine := scanner.Text()
		// ignore the first line (loadNumber pickup dropoff)
		if lineNum == 1 {
			continue
		}

		// contents -> [id, pickup, dropoff]
		contents := strings.Split(newLine, " ")

		// get id
		id, err := strconv.Atoi(contents[0])
		if err != nil {
			log.Fatal("Unable to convert ID to int in data: ", err)
		}

		// get pickup and dropoff coordinates
		pickUp := getCoord(contents[1])
		dropOff := getCoord(contents[2])

		// initialize a single load and store it to loads
		load := &LoadSingle{
			Pickup:  pickUp,
			Dropoff: dropOff,
		}
		loadsData[id] = load
	}

	return loadsData
}

// Preprocess the coordinates (can use regex as well)
func getCoord(coord string) *Coordinates {
	// remove brackets and seperate x and y
	coord = coord[1 : len(coord)-1]
	coordList := strings.Split(coord, ",")

	// convert x and y to float
	x, err := strconv.ParseFloat(coordList[0], 64)
	if err != nil {
		log.Fatal("Unable to convert x coordinate to float: ", err)
	}
	y, err := strconv.ParseFloat(coordList[1], 64)
	if err != nil {
		log.Fatal("Unable to convert y coordinate to float: ", err)
	}

	// Make and return coordinates
	coords := &Coordinates{
		X: x,
		Y: y,
	}

	return coords
}
