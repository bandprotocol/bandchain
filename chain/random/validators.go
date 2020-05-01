package main

import (
	"crypto/sha256"
	"encoding/binary"
)

type Validator struct {
	ID          []byte
	VotingPower int64
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
	luckyIdx := int64(binary.BigEndian.Uint64(seed)) % randomRange
	if luckyIdx <= 0 {
		luckyIdx *= int64(-1)
	}
	luckyVal := vals.ValidatorSet[int(luckyIdx)]
	vals.ValidatorSet = append(vals.ValidatorSet[:luckyIdx], vals.ValidatorSet[luckyIdx+1:]...)

	return luckyVal, vals
}

func getRandomRange(vals Validators, round, algo int) int64 {
	var randomRange int64
	switch algo {
	case 0:
		randomRange = int64(len(vals.ValidatorSet) - 10*round)
		if randomRange < 0 {
			randomRange = 0
		}
		if randomRange >= int64(len(vals.ValidatorSet)) {
			randomRange = int64(len(vals.ValidatorSet)) - 1
		}
	default:
		panic("ERROR")
	}
	return randomRange
}

func randomValidators(seed []byte, vals Validators, amount int, algo int) Validators {
	vals.GetVotingPower()

	var luckyVal Validators
	var val Validator
	for round := 0; round < amount; round++ {
		seed = nextSeed(seed)
		randomRange := getRandomRange(vals, round, algo)
		val, vals = luckyDraw(seed, vals, randomRange)
		luckyVal.ValidatorSet = append(luckyVal.ValidatorSet, val)
	}

	return luckyVal
}
