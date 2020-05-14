package bandrng

import "math"

// addUint64Overflow performs the addition operation on two uint64 integers and
// returns a boolean on whether or not the result overflows.
func addUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

// ChooseOne randomly picks an index between 0 and len(weights)-1 inclusively. Each index has
// the probability of getting selected based on its weight.
func ChooseOne(rng *Rng, weights []uint64) int {
	sum := uint64(0)
	var overflow bool
	for _, weight := range weights {
		sum, overflow = addUint64Overflow(sum, weight)
		if overflow {
			panic("sum of weights exceed max uint64")
		}
	}

	luckyNumber := rng.NextUint64() % sum
	currentSum := uint64(0)
	for idx, weight := range weights {
		currentSum += weight
		if currentSum > luckyNumber {
			return idx
		}
	}
	// Should never happen because the sum of weights is more than the lucky number
	panic("error")
}

// GetCandidateSize return candidate size that base on current round and total round
// currentRound must in range [0,totalRound)
// totalCount must be more than 0
// if totalRound = 1 the function will return totalCount
// candidate size will decrease every round
// candidate size calculate by function
// size = floor((totalCount - 1)**((totalRound-currentRound-1)/(totalRound-1))) + 1
// so size must in range [2,totalCount]
func GetCandidateSize(currentRound, totalRound, totalCount int) int {
	if currentRound < 0 || currentRound >= totalRound {
		panic("currentRound must in range [0,totalRound)")
	}
	if totalCount <= 0 {
		panic("totalCount must be more than 0")
	}

	if totalRound == 1 {
		return totalCount
	}

	size := math.Pow(float64(totalCount-1), float64(totalRound-currentRound-1)/float64(totalRound-1)) + float64(1)
	return int(size)
}
