package main

import (
	"github.com/Chuvvi/Vorto/algo"
	"github.com/Chuvvi/Vorto/preprocess"
)

func main() {
	// preprocess and get data
	loadsData := preprocess.GetData()

	algo.Greedy(loadsData)
}
