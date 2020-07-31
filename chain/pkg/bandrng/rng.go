package bandrng

import (
	"crypto"
	"encoding/binary"

	"github.com/oasisprotocol/oasis-core/go/common/crypto/drbg"
)

// Rng implements a simple determinisic random number generator. Starting from an initial entropy,
// nonce, and personalizationString, it utilizes HMAC_DRBG construct as per NIST Special
// Publication 800-90A to produce a stream of random uint64 integers.
type Rng struct {
	rng *drbg.Drbg
}

// NewRng creates a new psudo-random generator, using the given seeds as the initial random source.
func NewRng(entropyInput, nonce, personalizationString []byte) (*Rng, error) {
	rng, err := drbg.New(crypto.SHA256, entropyInput, nonce, personalizationString)
	if err != nil {
		return nil, err
	}
	return &Rng{rng: rng}, nil
}

// NextUint64 returns the next 64-bit unsigned random integer produced by this generator.
func (r *Rng) NextUint64() uint64 {
	data := make([]byte, 8)
	_, err := r.rng.Read(data)
	if err != nil {
		// Reaching error is not practically possible in hmbc_drbg's codepath.
		panic(err)
	}
	return binary.BigEndian.Uint64(data)
}
