package bandrng_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/pkg/bandrng"
)

func TestSamplingOne(t *testing.T) {
	r := bandrng.NewRng("SEED")
	weights := []uint64{10, 13, 10, 25, 42} // prefix sum is 10,25,55,80,100

	require.Equal(t, bandrng.SamplingOne(r, weights), 4) // rng NextUint64() will return 15735084640102210068
	require.Equal(t, bandrng.SamplingOne(r, weights), 4) // rng NextUint64() will return 3485776390957061973
	require.Equal(t, bandrng.SamplingOne(r, weights), 3) // rng NextUint64() will return 17609118114147816341
	require.Equal(t, bandrng.SamplingOne(r, weights), 2) // rng NextUint64() will return 15960811988050104523
	require.Equal(t, bandrng.SamplingOne(r, weights), 3) // rng NextUint64() will return 11919533627209787235

}

func TestSamplingOnePanic(t *testing.T) {
	r := bandrng.NewRng("SEED")

	require.Panics(t, func() {
		bandrng.SamplingOne(r, []uint64{math.MaxUint64, math.MaxUint64})
	})

	require.Panics(t, func() {
		bandrng.SamplingOne(r, []uint64{1, math.MaxUint64})
	})

	require.Panics(t, func() {
		bandrng.SamplingOne(r, []uint64{math.MaxUint64, 1})
	})

}

func TestaddUint64Overflow(t *testing.T) {

	sum, overflow := bandrng.AddUint64Overflow(1, 2)
	require.False(t, overflow)
	require.Equal(t, sum, uint64(3))

	sum, overflow = bandrng.AddUint64Overflow(math.MaxUint64, 0)
	require.False(t, overflow)
	require.Equal(t, sum, uint64(math.MaxUint64))

	sum, overflow = bandrng.AddUint64Overflow(0, math.MaxUint64)
	require.False(t, overflow)
	require.Equal(t, sum, uint64(math.MaxUint64))

	sum, overflow = bandrng.AddUint64Overflow(math.MaxUint64, math.MaxUint64)
	require.True(t, overflow)
	require.Equal(t, sum, uint64(0))

	sum, overflow = bandrng.AddUint64Overflow(math.MaxUint64, 1)
	require.True(t, overflow)
	require.Equal(t, sum, uint64(0))

	sum, overflow = bandrng.AddUint64Overflow(uint64(1), math.MaxUint64)
	require.True(t, overflow)
	require.Equal(t, sum, uint64(0))
}
