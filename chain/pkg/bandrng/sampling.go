package bandrng

import (
	"math"
)

// safeAdd performs the addition operation on two uint64 integers, but panics if overflow.
func safeAdd(a, b uint64) uint64 {
	if math.MaxUint64-a < b {
		panic("bandrng::safeAdd: overflow addition")
	}
	return a + b
}

// ChooseOne randomly picks an index between 0 and len(weights)-1 inclusively. Each index has
// the probability of getting selected based on its weight.
func ChooseOne(rng *Rng, weights []uint64) int {
	sum := uint64(0)
	for _, weight := range weights {
		sum = safeAdd(sum, weight)
	}

	luckyNumber := rng.NextUint64() % sum
	currentSum := uint64(0)
	for idx, weight := range weights {
		currentSum += weight
		if currentSum > luckyNumber {
			return idx
		}
	}
	// We should never reach here since the sum of weights is greater than the lucky number.
	panic("bandrng::ChooseOne: reaching the unreachable")
}

// ChooseSome randomly picks non-duplicate "cnt" indexes between 0 and len(weights)-1 inclusively.
// The function calls ChooseOne to get an index based on the given weights. When an index is
// chosen, it gets removed from the pool. The process gets repeated until "cnt" indexes are chosen.
func ChooseSome(rng *Rng, weights []uint64, cnt int) []int {
	chosenIndexes := make([]int, cnt)
	availableWeights := make([]uint64, len(weights))
	availableIndexes := make([]int, len(weights))
	for idx, weight := range weights {
		availableWeights[idx] = weight
		availableIndexes[idx] = idx
	}
	for round := 0; round < cnt; round++ {
		chosen := ChooseOne(rng, availableWeights)
		chosenIndexes[round] = availableIndexes[chosen]
		availableWeights = append(availableWeights[:chosen], availableWeights[chosen+1:]...)
		availableIndexes = append(availableIndexes[:chosen], availableIndexes[chosen+1:]...)
	}
	return chosenIndexes
}

// ChooseSomeMaxWeight performs ChooseSome "tries" times and returns the sampling with the
// highest weight sum among all tries.
func ChooseSomeMaxWeight(rng *Rng, weights []uint64, cnt int, tries int) []int {
	var maxWeightSum uint64 = 0
	var maxWeightResult []int = nil
	for each := 0; each < tries; each++ {
		candidate := ChooseSome(rng, weights, cnt)
		candidateWeightSum := uint64(0)
		for _, idx := range candidate {
			candidateWeightSum += weights[idx]
		}
		if candidateWeightSum > maxWeightSum {
			maxWeightSum = candidateWeightSum
			maxWeightResult = candidate
		}
	}
	return maxWeightResult
}
