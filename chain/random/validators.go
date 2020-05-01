package main

import (
	"crypto/sha256"
	"encoding/binary"
	"math/rand"
)

type Validator struct {
	Id          []byte
	votingPower int64
}

type Validators struct {
	validatorSet []Validator
}

// range specification, note that min <= max
type RandomRange struct {
	min, max int64
}

func (vals Validators) GetVotingPower() int64 {
	sum := int64(0)
	for _, val := range vals.validatorSet {
		sum += val.votingPower
	}
	return sum
}

func (vals Validators) Append(val Validator) {
	vals.validatorSet = append(vals.validatorSet, val)
}

func createRandomSeed(seedByte []byte) int64 {
	h := sha256.New()
	h.Write(seedByte)
	seed := h.Sum(nil)

	h = sha256.New()
	h.Write(seed[len(seed)-3:])
	seed = h.Sum(nil)

	h = sha256.New()
	h.Write(seed[len(seed)-3:])
	seed = h.Sum(nil)

	return int64(binary.BigEndian.Uint64(seed))
}
func (rr *RandomRange) NextRandom(r *rand.Rand) int64 {
	return r.Int63n(rr.max-rr.min+1) + rr.min
}

func luckyDraw(seedByte []byte, vals Validators, randomRange int64) (Validator, Validators) {
	r := rand.New(rand.NewSource(createRandomSeed(seedByte)))
	rr := RandomRange{0, randomRange}

	luckyIdx := rr.NextRandom(r)
	luckyVal := vals.validatorSet[int(luckyIdx)]
	vals.validatorSet = append(vals.validatorSet[:luckyIdx], vals.validatorSet[luckyIdx+1:]...)

	return luckyVal, vals
}

func getRandomRage(vals Validators, round, algo int) int64 {

	var randomRange int64
	switch algo {
	case 0:
		randomRange = int64(len(vals.validatorSet) - 10*round)
		if randomRange < 0 {
			randomRange = 0
		}
		if randomRange >= int64(len(vals.validatorSet)) {
			randomRange = int64(len(vals.validatorSet)) - 1
		}
	default:
		panic("ERROR")
	}
	return randomRange
}

func randomValidators(seedByte []byte, vals Validators, amount int, algo int) Validators {
	vals.GetVotingPower()

	var luckyVal Validators
	var val Validator
	for round := 0; round < amount; round++ {
		randomRange := getRandomRage(vals, round, algo)
		val, vals = luckyDraw(seedByte, vals, randomRange)
		luckyVal.validatorSet = append(luckyVal.validatorSet, val)
	}

	return luckyVal
}
