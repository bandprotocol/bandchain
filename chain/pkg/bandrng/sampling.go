package bandrng

import "math"

// AddUint64Overflow performs the addition operation on two uint64 integers and
// returns a boolean on whether or not the result overflows.
func AddUint64Overflow(a, b uint64) (uint64, bool) {
	if math.MaxUint64-a < b {
		return 0, true
	}

	return a + b, false
}

// SamplingOne sampling an index weighted by probability
func SamplingOne(rng *Rng, weights []uint64) int {
	sum := uint64(0)
	var overflow bool
	for _, weight := range weights {
		sum, overflow = AddUint64Overflow(sum, weight)
		if overflow {
			panic("a sum of weight is more than max uint64")
		}
	}

	luckyNumber := rng.NextUint64() % sum
	count := uint64(0)
	for idx, weight := range weights {
		count, overflow = AddUint64Overflow(count, weight)
		if overflow {
			panic("a um of weight is more than max uint64")

		}
		if count > luckyNumber {
			return idx
		}
	}
	// Should never happen because the sum of weight is more than a lucky number
	panic("error")
}
