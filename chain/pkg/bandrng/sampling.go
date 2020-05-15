package bandrng

import (
	"math"
)

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
// totalRound must be more than 0 and totalCount <= totalRound
// if currentRound is 0 the function will return totalCount
// candidate size will decrease every round
// candidate size calculate by function
// size = floor((totalCount-1)**((totalRound-currentRound-1)/(totalRound-1))) + 1
// so size must in range [2,totalCount]
func GetCandidateSize(currentRound, totalRound, totalCount int) int {
	if currentRound < 0 || currentRound >= totalRound {
		panic("currentRound must in range [0,totalRound)")
	}
	if totalCount < totalRound {
		panic("error: totalCount < totalRound")
	}

	if currentRound == 0 {
		return totalCount
	}

	base := float64(totalCount - 1)                                        // base > 0
	exponent := float64(totalRound-1-currentRound) / float64(totalRound-1) // 0 <= exponent <= 1

	size := int(math.Pow(base, exponent)) + 1
	return size
}

// ChooseK randomly picks an array of index(size=k) between 0 and len(weights)-1 inclusively. Each index has
// the probability of getting selected based on its weight.
func ChooseK(rng *Rng, weights []uint64, k int) []int {
	var luckies []int
	chooses := make([]bool, len(weights))
	ws := weights
	for round := 0; round < k; round++ {
		candidateSize := GetCandidateSize(round, k, len(ws))
		luckyNumber := ChooseOne(rng, ws[:candidateSize])
		if luckyNumber == 0 {
			ws = ws[luckyNumber+1:]
		} else {
			ws = append(ws[:luckyNumber], ws[luckyNumber+1:]...)
		}

		sum := 0
		for idx, choose := range chooses {
			if !choose {
				sum++
			}
			if sum == luckyNumber+1 {
				chooses[idx] = true
				luckies = append(luckies, idx)
				break
			}

		}
	}

	return luckies
}
