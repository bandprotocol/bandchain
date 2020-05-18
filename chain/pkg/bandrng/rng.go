package bandrng

import (
	"crypto/sha256"
	"encoding/binary"
)

type Rng struct {
	seed [sha256.Size]byte
}

func nextSeed(prev [sha256.Size]byte) [sha256.Size]byte {
	return sha256.Sum256(prev[:])
}

// NewRng creates a new psudo-random generator, using the given seed as the initial random source.
func NewRng(initSeed string) *Rng {
	return &Rng{seed: sha256.Sum256([]byte(initSeed))}
}

// NextUint64 returns the next 64-bit unsigned random integer produced by this generator.
func (r *Rng) NextUint64() uint64 {
	val := binary.BigEndian.Uint64(r.seed[:8])
	r.seed = nextSeed(r.seed)
	return val
}
