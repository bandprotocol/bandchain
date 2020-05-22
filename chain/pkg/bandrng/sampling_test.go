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

func TestGetCandidateSize(t *testing.T) {
	totalRound := 7
	require.Equal(t, bandrng.GetCandidateSize(0, totalRound, 93), 93)
	require.Equal(t, bandrng.GetCandidateSize(1, totalRound, 92), 43)
	require.Equal(t, bandrng.GetCandidateSize(2, totalRound, 91), 21)
	require.Equal(t, bandrng.GetCandidateSize(3, totalRound, 90), 10)
	require.Equal(t, bandrng.GetCandidateSize(4, totalRound, 89), 5)
	require.Equal(t, bandrng.GetCandidateSize(5, totalRound, 88), 3)
	require.Equal(t, bandrng.GetCandidateSize(6, totalRound, 87), 2)

	require.Equal(t, bandrng.GetCandidateSize(0, 1, 9999), 9999)
	require.Equal(t, bandrng.GetCandidateSize(1, 2, 9999), 2)

}

func TestGetCandidateSizePanic(t *testing.T) {
	require.Panics(t, func() {
		bandrng.GetCandidateSize(10, 0, 99)
	})
	require.Panics(t, func() {
		bandrng.GetCandidateSize(10, 10, 99)
	})
	require.Panics(t, func() {
		bandrng.GetCandidateSize(9, 10, 0)
	})
}

func TestChooseK(t *testing.T) {
	r := bandrng.NewRng("SEED")
	length := 93
	weights := make([]uint64, length)
	for idx := 0; idx < length; idx++ {
		weights[idx] = 1
	}

	acc := bandrng.ChooseK(r, weights, 7)
	require.Equal(t, acc, []int{84, 13, 8, 3, 0, 2, 4})

}
