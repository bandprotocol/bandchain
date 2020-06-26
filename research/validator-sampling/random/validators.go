package main

import (
	"crypto/sha256"
	"encoding/binary"
	"math"
)

type Validator struct {
	ID          []byte
	VotingPower int64
	Number      int
}

type Validators struct {
	ValidatorSet []Validator
}

func (vals Validators) GetVotingPower() int64 {
	sum := int64(0)
	for _, val := range vals.ValidatorSet {
		sum += val.VotingPower
	}
	return sum
}

func nextSeed(seed []byte) []byte {
	h := sha256.New()
	h.Write(seed)
	seed = h.Sum(nil)

	h = sha256.New()
	h.Write(seed[len(seed)-3:])
	seed = h.Sum(nil)

	h = sha256.New()
	h.Write(seed[len(seed)-3:])
	seed = h.Sum(nil)

	return seed
}

func luckyDraw(seed []byte, vals Validators, randomRange int64) (Validator, Validators) {
	luckyNumber := int64(binary.BigEndian.Uint64(seed)) % randomRange
	if luckyNumber <= 0 {
		luckyNumber *= int64(-1)
	}
	sum := int64(0)
	luckyIdx := 0
	for idx := 0; sum < luckyNumber && idx < len(vals.ValidatorSet); idx++ {
		sum += vals.ValidatorSet[idx].VotingPower
		luckyIdx = idx
	}
	luckyVal := vals.ValidatorSet[luckyIdx]
	if luckyIdx != 0 {
		vals.ValidatorSet = append(vals.ValidatorSet[:luckyIdx], vals.ValidatorSet[luckyIdx+1:]...)
	} else {
		vals.ValidatorSet = vals.ValidatorSet[luckyIdx+1:]
	}
	return luckyVal, vals
}

func getRandomRange(vals Validators, round, amount, num, algo int) int64 {
	randomRange := int64(0)
	switch algo {

	case 0:
		blockSize := num / amount

		t := (amount - round) * blockSize
		if round == 0 {
			t = len(vals.ValidatorSet)
		}

		for idx := 0; idx < t; idx++ {
			randomRange += vals.ValidatorSet[idx].VotingPower
		}
		if randomRange == 0 {
			randomRange = 1
		}
	case 1:
		t := num / (1 << round)

		for idx := 0; idx < t; idx++ {
			randomRange += vals.ValidatorSet[idx].VotingPower
		}
		if randomRange == 0 {
			randomRange = 1
		}
	case 2:
		c := 5
		n := len(vals.ValidatorSet)
		m := float64(n-c) / float64(amount-1)
		t := int(m*float64(amount-round-1) + float64(c))

		// fmt.Println("->", t)
		for idx := 0; idx < t; idx++ {
			randomRange += vals.ValidatorSet[idx].VotingPower
		}
		if randomRange == 0 {
			randomRange = 1
		}
	case 3:
		c := amount - round + 1
		// c := 2
		n := len(vals.ValidatorSet)
		x := float64(amount-round-1) / float64(amount-1) * math.Log(float64(n-c))
		t := math.Exp(x) + float64(c-1)
		if round == 0 {
			t = float64(n)
		}
		if t > float64(n) {
			t = float64(n)
		}
		// fmt.Println("->", t)

		for idx := 0; idx < int(t); idx++ {
			randomRange += vals.ValidatorSet[idx].VotingPower
		}
		if randomRange == 0 {
			randomRange = 1
		}
	case 4:
		c := amount - round + 1
		// c := 2
		n := len(vals.ValidatorSet)
		x := float64(amount-round-1) / float64(amount-1) * math.Log(float64(n-c))
		t := math.Exp(x) + float64(c-1)
		if round == 0 {
			t = float64(n)
		}
		if t > float64(n) {
			t = float64(n)
		}
		// fmt.Println("->", t)

		for idx := 0; idx < int(t); idx++ {
			randomRange += vals.ValidatorSet[idx].VotingPower
		}
		if randomRange == 0 {
			randomRange = 1
		}
	case 5:
		// c := amount - round + 1
		c := 2
		n := len(vals.ValidatorSet)
		x := float64(amount-round-1) / float64(amount-1) * math.Log(float64(n-c))
		t := math.Exp(x) + float64(c-1)
		if round == 0 {
			t = float64(n)
		}
		if t > float64(n) {
			t = float64(n)
		}
		// fmt.Println("->", t)

		for idx := 0; idx < int(t); idx++ {
			randomRange += vals.ValidatorSet[idx].VotingPower
		}
		if randomRange == 0 {
			randomRange = 1
		}
	case 6:
		// blockSize := num / amount

		t := (num - amount)
		if round == 0 {
			t = len(vals.ValidatorSet)
		}

		for idx := 0; idx < t; idx++ {
			randomRange += vals.ValidatorSet[idx].VotingPower
		}

	default:
		panic("ERROR")
	}

	return randomRange
}

func getWorseCaseVp(vals Validators, amount, num, algo int) int64 {
	sumVotingPower := int64(0)

	switch algo {
	case 0:
		blockSize := len(vals.ValidatorSet) / amount

		for round := 0; round < amount; round++ {
			t := (round + 1) * blockSize
			sumVotingPower += vals.ValidatorSet[t-1].VotingPower
		}
	case 1:
		ct := 0
		use := make(map[int]bool)
		for round := 0; round < amount; round++ {
			t := num / (1 << round)
			if t-1 > 0 {
				sumVotingPower += vals.ValidatorSet[t-1].VotingPower
				ct++
				use[t-1] = true
			}
		}
		for idx := 0; ct < amount; idx++ {
			if _, ok := use[idx]; !ok {
				sumVotingPower += vals.ValidatorSet[idx].VotingPower
				ct++
				use[idx] = true
			}
		}
	case 2:
		for round := 0; round < amount; round++ {
			c := 5
			n := len(vals.ValidatorSet)
			m := float64(n-c) / float64(amount-1)
			t := int(m*float64(amount-round-1) + float64(c))
			// fmt.Println("!!->", t, m)
			if t > n {
				t = n
			}
			sumVotingPower += vals.ValidatorSet[t-1].VotingPower
		}
	case 3:

		for round := 0; round < amount; round++ {
			c := amount - round + 1
			// c := 2
			n := len(vals.ValidatorSet)
			x := float64(amount-round-1) / float64(amount-1) * math.Log(float64(n-c))
			t := math.Exp(x) + float64(c-1)

			if t > float64(n) {
				t = float64(n)
			}
			// fmt.Println("->", t)

			sumVotingPower += vals.ValidatorSet[int(t)-1].VotingPower
		}
	case 4:

		for round := 0; round < amount; round++ {
			c := amount - round + 1
			// c := 2
			n := len(vals.ValidatorSet)
			x := float64(amount-round-1) / float64(amount-1) * math.Log(float64(n-c))
			t := math.Exp(x) + float64(c-1)

			if t > float64(n) {
				t = float64(n)
			}
			// fmt.Println("->", t)

			sumVotingPower += vals.ValidatorSet[int(t)-1].VotingPower
		}
	case 5:

		for round := 0; round < amount; round++ {
			// c := amount - round + 1
			c := 2
			n := len(vals.ValidatorSet)
			x := float64(amount-round-1) / float64(amount-1) * math.Log(float64(n-c))
			t := math.Exp(x) + float64(c-1)

			if t > float64(n) {
				t = float64(n)
			}
			// fmt.Println("->", t)

			sumVotingPower += vals.ValidatorSet[int(t)-1].VotingPower
		}
	case 6:
		sumVotingPower = 0
	default:
		panic("ERROR")
	}
	return sumVotingPower
}

func randomValidators(seed []byte, vals Validators, amount int, algo int) (Validators, int64) {
	vals.GetVotingPower()

	var luckyVal Validators
	var val Validator
	num := len(vals.ValidatorSet)

	// topBadCase := make([]int64, 0)

	worseCase := getWorseCaseVp(vals, amount, num, algo)
	for round := 0; round < amount; round++ {
		seed = nextSeed(seed)
		randomRange := getRandomRange(vals, round, amount, num, algo)
		val, vals = luckyDraw(seed, vals, randomRange)
		luckyVal.ValidatorSet = append(luckyVal.ValidatorSet, val)

	}

	// topBadCase = append(topBadCase, luckyVal.GetVotingPower())
	// fmt.Println(topBadCase)
	// if len(topBadCase) == top {
	// 	sort.Slice(topBadCase, func(i, j int) bool { return topBadCase[i] < topBadCase[j] })
	// 	topBadCase = topBadCase[:top]
	// }

	return luckyVal, worseCase
}
