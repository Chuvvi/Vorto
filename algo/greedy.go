package algo

import (
	"fmt"

	"github.com/Chuvvi/Vorto/preprocess"
)

var distArr [][]float64

func Greedy(data preprocess.Loads) {
	drivers1, c1 := GreedyDriverToDelivery(data)
	drivers2, c2 := GreedyDeliveryToDriver(data)
	if c1 < c2 {
		PrintOP(drivers1)
	} else {
		PrintOP(drivers2)
	}
}

// function to calculate the cost of all drivers
func calculateCost(data preprocess.Loads, drivers map[int]*Driver) float64 {
	cost := 0.0
	for _, driver := range drivers {
		currPos := 0
		for _, delivery := range driver.Deliveries {
			if currPos == 0 {
				cost += getTime(*DEPOT_COORD, *data[delivery])
			} else {
				cost += getTime(*data[currPos].Dropoff, *data[delivery])
			}
			currPos = delivery
		}
		cost += euclidianDistance(*data[currPos].Dropoff, *DEPOT_COORD)
		// cost += driver.Time
	}
	return cost
}

// populate distArr with distances between each deliveries and/or depot
func initializeDist(data preprocess.Loads) {
	n := len(data) + 1
	for i := 0; i < n; i++ {
		distArr[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			if i == j {
				distArr[i][j] = 0.0
			} else if i == 0 {
				distArr[i][j] = getTime(*DEPOT_COORD, *data[j])
			} else if j == 0 {
				distArr[i][j] = euclidianDistance(*data[i].Dropoff, *DEPOT_COORD)
			} else {
				distArr[i][j] = getTime(*data[i].Dropoff, *data[j])
			}
		}
	}
}

// driver function to print out the output
func PrintOP(drivers map[int]*Driver) {
	for _, driverData := range drivers {
		fmt.Printf("[")
		n := len(driverData.Deliveries)
		for i := 0; i < n-1; i++ {
			fmt.Printf("%v,", driverData.Deliveries[i])
		}
		fmt.Printf("%v]\n", driverData.Deliveries[n-1])
	}
}
