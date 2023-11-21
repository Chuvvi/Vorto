package algo

import (
	"math/rand"
	"sort"
	"vorto/preprocess"
)

type Indices struct {
	Idx     int
	RandNum float64
}

type ChromoCost struct {
	Cost   float64
	Chromo []*Indices
}

// implements the genetic algorithm
func Genetic(data preprocess.Loads, k int, iterations int, hp float64) (float64, []int) {
	// generate and store k chromosomes

	chromosomes := make([][]*Indices, 0)
	for i := 0; i < k; i++ {
		chromo := newChromo(data)
		chromosomes = append(chromosomes, chromo)
	}

	// Iterate "iterations" times
	for iter := 0; iter < iterations; iter++ {
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
			child := newChromo(data)
			newChildren = append(newChildren, child)
		}

		// now we have 2*k chromosomes
		// select the best k chromosomes from this
		chromoOp := make([]*ChromoCost, 0)
		chromoOp = append(chromoOp, evalChromo(data, chromosomes, hp)...)
		chromoOp = append(chromoOp, evalChromo(data, crossOverChildren, hp)...)
		chromoOp = append(chromoOp, evalChromo(data, mutationChildren, hp)...)
		chromoOp = append(chromoOp, evalChromo(data, newChildren, hp)...)
		sort.Slice(chromoOp, func(i, j int) bool {
			return chromoOp[i].Cost < chromoOp[j].Cost
		})

		chromosomes = make([][]*Indices, 0)
		for i := 0; i < k; i++ {
			chromosomes = append(chromosomes, chromoOp[k].Chromo)
		}
	}

	// get the best chromosome and display the output
	bestChromo := chromosomes[0]
	loadIDs := make([]int, 0)
	for _, v := range bestChromo {
		loadIDs = append(loadIDs, v.Idx)
	}
	cost := greedy(data, loadIDs, hp, false)
	return cost, loadIDs
}

// generate a random new chromosome
func newChromo(data preprocess.Loads) []*Indices {
	chromo := make([]*Indices, 0)
	for k := range data {
		chromo = append(chromo, &Indices{
			Idx:     k,
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

// driver function to estimate the cost of each chromosome
func evalChromo(data preprocess.Loads, chromosomes [][]*Indices, hp float64) []*ChromoCost {
	chromoOP := make([]*ChromoCost, 0)
	for _, chromo := range chromosomes {
		loadIDs := make([]int, 0)
		for _, v := range chromo {
			loadIDs = append(loadIDs, v.Idx)
		}
		cost := greedy(data, loadIDs, hp, false)
		chromoOP = append(chromoOP, &ChromoCost{
			Cost:   cost,
			Chromo: chromo,
		})
	}
	return chromoOP
}
