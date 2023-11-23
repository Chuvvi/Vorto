package algo

import (
	"math"

	"github.com/Chuvvi/Vorto/preprocess"
)

// assign delivery to driver
func GreedyDeliveryToDriver(data preprocess.Loads) (map[int]*Driver, float64) {
	// make a adjacency list for easier calculation
	distArr = make(([][]float64), len(data)+1)
	initializeDist(data)
	bestDrivers := make(map[int]*Driver)
	c := math.MaxFloat64
	bestDriversArr := assignDeliveries(data)
	for _, bd := range bestDriversArr {
		// more optimization can be done at the cost of processing power
		// if genetic algorithm is removed execution will finish in about 30ms
		bd := Genetic(data, bd, 50, 3)
		cost := calculateCost(data, bd)
		if cost < c {
			c = cost
			bestDrivers = bd
		}
	}
	return bestDrivers, c
}

func assignDeliveries(data preprocess.Loads) []map[int]*Driver {
	// make a slice of drivers starting from depot and each of them coonecting to a delivery
	driversArr := make([]map[int]*Driver, 0)
	n := len(data)
	// start assigning deliveries
	for i := 1; i <= n; i++ {
		// keep a track of deliveries
		deliveriesMade := make(map[int]bool)
		// create a new driver
		drivers := make(map[int]*Driver)
		driver := &Driver{
			CurrCoord:  *data[i].Dropoff,
			Time:       distArr[0][i],
			Deliveries: make([]int, 0),
		}
		driver.Deliveries = append(driver.Deliveries, i)
		deliveriesMade[i] = true
		currDeliveryID := i
		dID := 0
		// deliver until all deliveries are made
		for len(deliveriesMade) < n {
			for {
				// find the next closest delivery
				minDist := math.MaxFloat64
				minLoadID := -1
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
				// if minDist hasn't changed then no new deliveries can be made
				if minDist == math.MaxFloat64 {
					break
				}
				// update driver
				currDeliveryID = minLoadID
				driver.CurrCoord = *data[minLoadID].Dropoff
				driver.Time += minDist
				driver.Deliveries = append(driver.Deliveries, currDeliveryID)
				deliveriesMade[currDeliveryID] = true
			}
			// send the driver back to depot
			driver.Time += euclidianDistance(driver.CurrCoord, *DEPOT_COORD)
			driver.CurrCoord = *DEPOT_COORD
			drivers[dID] = driver
			dID++
			// create a new driver starting from the depot
			driver = &Driver{
				CurrCoord:  *DEPOT_COORD,
				Time:       0,
				Deliveries: make([]int, 0),
			}
			currDeliveryID = 0
		}
		driversArr = append(driversArr, drivers)
	}
	return driversArr
}
