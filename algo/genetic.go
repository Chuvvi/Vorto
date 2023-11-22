package algo

import (
	"math/rand"
	"sort"

	"github.com/Chuvvi/Vorto/preprocess"
)

type Indices struct {
	Idx     int
	RandNum float64
}

type ChromoCost struct {
	Cost   float64
	Chromo []*Indices
}

// driver function to get total distance travelled by a driver
func getDist(data preprocess.Loads, deliveries []int) (float64, bool) {
	dist := 0.0
	currCoord := DEPOT_COORD
	for _, v := range deliveries {
		dist += euclidianDistance(*currCoord, *data[v].Pickup)
		dist += euclidianDistance(*data[v].Pickup, *data[v].Dropoff)
		currCoord = data[v].Dropoff
	}
	dist += euclidianDistance(*currCoord, *DEPOT_COORD)
	if dist > MAX_TIME {
		return -1, false
	}
	return dist, true
}

// implements the genetic algorithm
func Genetic(data preprocess.Loads, drivers map[int]*Driver, k int, generations int) map[int]*Driver {
	optimalDrivers := make(map[int]*Driver)
	for driverID, driver := range drivers {
		// no optimization can be done if the driver delivers a single load
		if len(driver.Deliveries) == 1 {
			optimalDrivers[driverID] = driver
			continue
		}
		minDelivery := driver.Deliveries
		minDist, _ := getDist(data, minDelivery)

		// generate and store k chromosomes
		chromosomes := make([][]*Indices, 0)
		for i := 0; i < k; i++ {
			chromo := newChromo(minDelivery)
			chromosomes = append(chromosomes, chromo)
		}

		// Iterate "generations" times
		for iter := 0; iter < generations; iter++ {
			// Split k such that: 40% -> CrossOver, 40% -> Mutation, 20% -> New Chromosomes
			crossOverCount := int(0.4 * float64(k))
			mutationCount := int(0.4 * float64(k))
			newChromoCount := k - (crossOverCount + mutationCount)

			// make crossovers
			crossOverChildren := make([][]*Indices, 0)
			for i := 0; i < crossOverCount; i++ {
				// select two parents randomly for crossover
				p1 := rand.Intn(len(chromosomes))
				p2 := rand.Intn(len(chromosomes))
				child := crossOver(chromosomes[p1], chromosomes[p2])
				crossOverChildren = append(crossOverChildren, child)
			}

			// make mutations
			mutationChildren := make([][]*Indices, 0)
			for i := 0; i < mutationCount; i++ {
				// select one parent randomly for mutation
				p := rand.Intn(len(chromosomes))
				child := mutation(chromosomes[p])
				mutationChildren = append(mutationChildren, child)
			}

			// make new mutations
			newChildren := make([][]*Indices, 0)
			for i := 0; i < newChromoCount; i++ {
				child := newChromo(minDelivery)
				newChildren = append(newChildren, child)
			}

			// now we have 2*k chromosomes
			// select the best k chromosomes from this
			chromoOP := make([]*ChromoCost, 0)
			chromoOP = evalChromo(chromosomes, chromoOP, data)
			chromoOP = evalChromo(crossOverChildren, chromoOP, data)
			chromoOP = evalChromo(mutationChildren, chromoOP, data)
			chromoOP = evalChromo(newChildren, chromoOP, data)
			if len(chromoOP) == 0 {
				driver.Deliveries = minDelivery
				continue
			}

			// pad chromoOP to k incase there aren't enough data sets
			chromoOPP := make([]*ChromoCost, len(chromoOP))
			copy(chromoOPP, chromoOP)
			for len(chromoOPP) < 2*k {
				for _, v := range chromoOP {
					chromoOPP = append(chromoOPP, v)
					if len(chromoOP) == 2*k {
						break
					}
				}
			}
			chromoOP = chromoOPP

			// get the top k chromosomes
			sort.Slice(chromoOP, func(i, j int) bool {
				return chromoOP[i].Cost < chromoOP[j].Cost
			})
			chromosomes = make([][]*Indices, 0)
			for i := 0; i < k; i++ {
				chromosomes = append(chromosomes, chromoOP[i].Chromo)
			}

			// check if mincost is lesser and update
			if chromoOP[0].Cost < minDist {
				minDist = chromoOP[0].Cost
				newMinDelivery := make([]int, 0)
				for _, d := range chromoOP[0].Chromo {
					newMinDelivery = append(newMinDelivery, d.Idx)
				}
				minDelivery = newMinDelivery
			}

			driver.Deliveries = minDelivery
		}
		optimalDrivers[driverID] = driver
	}
	return optimalDrivers
}

// evaluate chromosomes
func evalChromo(chromosomes [][]*Indices, chromoOP []*ChromoCost, data preprocess.Loads) []*ChromoCost {
	for _, chromo := range chromosomes {
		deliveries := make([]int, 0)
		for _, v := range chromo {
			deliveries = append(deliveries, v.Idx)
		}
		if dist, ok := getDist(data, deliveries); ok {
			chromoOP = append(chromoOP, &ChromoCost{
				Cost:   dist,
				Chromo: chromo,
			})
		}
	}
	return chromoOP
}

// generate a random new chromosome
func newChromo(minDelivery []int) []*Indices {
	chromo := make([]*Indices, 0)
	for _, v := range minDelivery {
		chromo = append(chromo, &Indices{
			Idx:     v,
			RandNum: rand.Float64(),
		})
	}
	// sort in ascending order based on random numbers
	sort.Slice(chromo, func(i, j int) bool {
		return chromo[i].RandNum < chromo[j].RandNum
	})
	return chromo
}

// create a new chromosome from two parent chromosome
func crossOver(p1, p2 []*Indices) []*Indices {
	// select k indices from p1 randomly
	idx1, idx2 := getRandomRange(len(p1))
	// store all those indices in hashmap for constant access
	indicesHashMap := make(map[int]int)
	for idx1 <= idx2 {
		idx := p1[idx1].Idx
		indicesHashMap[idx] = idx1
		idx1++
	}

	// make the crossover child
	child := make([]*Indices, 0)
	for _, v := range p2 {
		idx := v.Idx
		if _, ok := indicesHashMap[idx]; !ok {
			child = append(child, v)
		}
	}
	for _, v := range indicesHashMap {
		child = append(child, p1[v])
	}
	return child
}

// create a new chromosome from a single parent chromosome
func mutation(p []*Indices) []*Indices {
	// select k indices from p randomly
	idx1, idx2 := getRandomRange(len(p))
	child := make([]*Indices, len(p))
	copy(child, p)
	for idx1 < idx2 {
		child[idx1], child[idx2] = child[idx2], child[idx1]
		idx1++
		idx2--
	}
	return child
}

// driver function to get random range within k
func getRandomRange(k int) (int, int) {
	idx1 := rand.Intn(k)
	idx2 := rand.Intn(k)
	for idx1 == idx2 {
		idx1 = rand.Intn(k)
		idx2 = rand.Intn(k)
	}
	if idx1 > idx2 {
		idx1, idx2 = idx2, idx1
	}
	return idx1, idx2
}
