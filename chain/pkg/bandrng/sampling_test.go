package bandrng_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/pkg/bandrng"
)

func TestChooseOneOne(t *testing.T) {
	r := bandrng.NewRng("SEED")
	weights := []uint64{10, 13, 10, 25, 42} // prefix sum is 10,23,33,58,100

	require.Equal(t, bandrng.ChooseOne(r, weights), 4) // rng NextUint64() will return 15735084640102210068
	require.Equal(t, bandrng.ChooseOne(r, weights), 4) // rng NextUint64() will return 3485776390957061973
	require.Equal(t, bandrng.ChooseOne(r, weights), 3) // rng NextUint64() will return 17609118114147816341
	require.Equal(t, bandrng.ChooseOne(r, weights), 2) // rng NextUint64() will return 15960811988050104523
	require.Equal(t, bandrng.ChooseOne(r, weights), 3) // rng NextUint64() will return 11919533627209787235

	r = bandrng.NewRng("SEED")
	weights = []uint64{2, 4, 4} // prefix sum is 2,6,10

	require.Equal(t, bandrng.ChooseOne(r, weights), 2) // rng NextUint64() will return 15735084640102210068
	require.Equal(t, bandrng.ChooseOne(r, weights), 1) // rng NextUint64() will return 3485776390957061973
	require.Equal(t, bandrng.ChooseOne(r, weights), 0) // rng NextUint64() will return 17609118114147816341

}

func TestChooseOnePanic(t *testing.T) {
	r := bandrng.NewRng("SEED")

	require.Panics(t, func() {
		bandrng.ChooseOne(r, []uint64{math.MaxUint64, math.MaxUint64})
	})

	require.Panics(t, func() {
		bandrng.ChooseOne(r, []uint64{1, math.MaxUint64})
	})

	require.Panics(t, func() {
		bandrng.ChooseOne(r, []uint64{math.MaxUint64, 1})
	})

}
