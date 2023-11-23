package algo

import (
	"math"

	"github.com/Chuvvi/Vorto/preprocess"
)

// assign delivery to driver
func GreedyDeliveryToDriver(data preprocess.Loads) (map[int]*Driver, float64) {
	a := assignDeliveries(data)
	b := Genetic(data, a, 100, 20)
	c := calculateCost(b)
	return a, c
}

func assignDeliveries(data preprocess.Loads) map[int]*Driver {
	// make drivers map
	drivers := make(map[int]*Driver)
	n := len(data)
	// keep track of deliveries
	deliveriesMade := make(map[int]bool)
	for i := 0; i <= n; i++ {
		// assign a new driver
		driver := &Driver{
			CurrCoord:  *DEPOT_COORD,
			Time:       0.0,
			Deliveries: make([]int, 0),
		}

		// start delivering
		currDeliveryID := 0
		for {
			// set minDist to max
			minDist := math.MaxFloat64
			minLoadID := -1
			// check for the nearest delivery point
			for j := 1; j <= n; j++ {
				_, ok := deliveriesMade[j]
				dist := distArr[currDeliveryID][j]
				depotDist := distArr[j][0]
				if j == currDeliveryID || ok || driver.Time+dist+depotDist > MAX_TIME {
					continue
				}
				if dist < minDist {
					minDist = dist
					minLoadID = j
				}
			}

			// if minDist is still max, then no new delivery can be made
			if minDist == math.MaxFloat64 {
				break
			}
			// update the driver values accordingly
			currDeliveryID = minLoadID
			driver.CurrCoord = *data[minLoadID].Dropoff
			driver.Time += minDist
			driver.Deliveries = append(driver.Deliveries, currDeliveryID)
			deliveriesMade[currDeliveryID] = true
		}

		// send the driver back to the depot
		driver.Time += euclidianDistance(driver.CurrCoord, *DEPOT_COORD)
		driver.CurrCoord = *DEPOT_COORD
		drivers[i] = driver

		// if all deliveries are made, exit the loop
		if len(deliveriesMade) == n {
			break
		}
	}
	return drivers
}
