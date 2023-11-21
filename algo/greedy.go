package algo

import (
	"fmt"
	"vorto/preprocess"
)

func Greedy(data preprocess.Loads) {
	for k, v := range data {
		fmt.Printf("%v: ((%v), (%v)) \n", k, v.Pickup, v.Dropoff)
	}

	// create drivers
	// currDriverID := 0
	// drivers := make(map[int][]int)

	// deliveriesLeft = len(data)
}
