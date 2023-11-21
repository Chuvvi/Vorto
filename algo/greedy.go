package algo

import (
	"fmt"
	"math"
	"sort"
	"vorto/preprocess"
)

type Driver struct {
	CurrCoord  preprocess.Coordinates
	Time       float64
	Deliveries []int
}

// Depot coordinates
var DEPOT_COORD *preprocess.Coordinates = &preprocess.Coordinates{
	X: 0.0,
	Y: 0.0,
}

// Maximum allowed time
const MAX_TIME float64 = 12 * 60

type LoadIdDist struct {
	Id   int
	Dist float64
}

func Greedy(data preprocess.Loads) {
	// iterate through the map is ascending order of their distance from depot for consistancy by iterating over sorted keys
	loadIDs := make([]LoadIdDist, 0)
	for loadID, loadData := range data {
		getTime(*DEPOT_COORD, *loadData)
		loadIDs = append(loadIDs, LoadIdDist{
			Id:   loadID,
			Dist: getTime(*DEPOT_COORD, *loadData),
		})
	}
	sort.Slice(loadIDs, func(i, j int) bool {
		return loadIDs[i].Dist < loadIDs[j].Dist
	})

	// get the best hyperparameter
	cost := math.MaxFloat64
	hp := -1.0
	for i := 0.4; i <= 1.6; i += 0.01 {
		c := greedy(data, loadIDs, i, false)
		if c < cost {
			cost = c
			hp = i
		}
	}
	greedy(data, loadIDs, hp, true)
}

func greedy(data preprocess.Loads, loadIDs []LoadIdDist, hp float64, displayOP bool) float64 {
	// create drivers
	currDriverID := 0
	drivers := make(map[int]*Driver)

	// Check for available load
	for _, loadData := range loadIDs {
		loadID := loadData.Id
		load := data[loadID]
		driverID := -1

		// assume that the max time required to deliver a load is from the depot
		requiredTime := hp * getTime(*DEPOT_COORD, *load)

		// find a suitable driver
		for dID, driverData := range drivers {
			deliveryTime := getTime(driverData.CurrCoord, *load)
			dropoffToDepotTime := euclidianDistance(*load.Dropoff, *DEPOT_COORD)
			totalTime := deliveryTime + dropoffToDepotTime
			if driverData.Time+totalTime > MAX_TIME || deliveryTime > requiredTime {
				continue
			}

			// update minimum time required to transit and driverID
			requiredTime = deliveryTime
			driverID = dID
		}

		// if suitable driver was not found, create a new driver
		if driverID == -1 {
			driverID = currDriverID
			addDriver(&drivers, driverID)
			requiredTime = getTime(*DEPOT_COORD, *load)
			currDriverID++
		}

		// update driver details
		drivers[driverID].CurrCoord = *load.Dropoff
		drivers[driverID].Time += requiredTime
		drivers[driverID].Deliveries = append(drivers[driverID].Deliveries, loadID)
	}

	cost := float64(500 * len(drivers))
	cost = 0.0
	for _, driverData := range drivers {
		driverData.Time += euclidianDistance(driverData.CurrCoord, *DEPOT_COORD)
		cost += driverData.Time
	}

	if displayOP {
		printOP(drivers)
	}

	return cost
}

// driver function to add a new driver starting at (0, 0)
func addDriver(drivers *map[int]*Driver, id int) {
	(*drivers)[id] = &Driver{
		CurrCoord: preprocess.Coordinates{
			X: 0.0,
			Y: 0.0,
		},
		Time:       0,
		Deliveries: make([]int, 0),
	}
}

// driver function to get the time elapsed for delivery
func getTime(driver preprocess.Coordinates, load preprocess.LoadSingle) float64 {
	currToPickup := euclidianDistance(driver, *load.Pickup)
	pickupToDropoff := euclidianDistance(*load.Pickup, *load.Dropoff)
	return currToPickup + pickupToDropoff
}

// euclidian distance formula
func euclidianDistance(p1, p2 preprocess.Coordinates) float64 {
	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

// driver function to print out the output
func printOP(drivers map[int]*Driver) {
	for _, driverData := range drivers {
		fmt.Printf("[")
		n := len(driverData.Deliveries)
		for i := 0; i < n-1; i++ {
			fmt.Printf("%v,", driverData.Deliveries[i])
		}
		fmt.Printf("%v]\n", driverData.Deliveries[n-1])
	}
}
