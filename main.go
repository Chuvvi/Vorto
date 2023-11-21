package main

import (
	"vorto/algo"
	"vorto/preprocess"
)

func main() {
	// preprocess and get data
	loadsData := preprocess.GetData()

	algo.Greedy(loadsData)
	// algo.Genetic(loadsData, 100, 10)
}
