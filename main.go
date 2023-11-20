package main

import (
	"fmt"
	"vorto/preprocess"
)

func main() {
	// preprocess and get data
	loadsData := preprocess.GetData()

	for id, load := range loadsData {
		fmt.Println(id, load.Pickup, load.Dropoff)
	}
}
