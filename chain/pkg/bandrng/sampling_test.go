package bandrng_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/pkg/bandrng"
)

func TestChooseOneOne(t *testing.T) {
	r, err := bandrng.NewRng([]byte("THIS_IS_A_RANDOM_SEED_LONG_ENOUGH_FOR_ENTROPY"), []byte("1"), []byte("TEST"))
	require.NoError(t, err)
	weights := []uint64{10, 13, 10, 25, 42} // prefix sum is 10,23,33,58,100

	require.Equal(t, 4, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 5751621621077249396
	require.Equal(t, 4, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 16474548556352052882
	require.Equal(t, 1, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 17097048712898369316
	require.Equal(t, 2, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 10686498023352306525
	require.Equal(t, 4, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 2144097648649487685

	r, err = bandrng.NewRng([]byte("THIS_IS_A_RANDOM_SEED_LONG_ENOUGH_FOR_ENTROPY"), []byte("1"), []byte("TEST"))
	require.NoError(t, err)
	weights = []uint64{2, 4, 4} // prefix sum is 2,6,10

	require.Equal(t, 2, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 5751621621077249396
	require.Equal(t, 1, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 16474548556352052882
	require.Equal(t, 2, bandrng.ChooseOne(r, weights)) // rng NextUint64() will return 17097048712898369316

}

func TestChooseOnePanic(t *testing.T) {
	r, err := bandrng.NewRng([]byte("THIS_IS_A_RANDOM_SEED_LONG_ENOUGH_FOR_ENTROPY"), []byte("1"), []byte("TEST"))
	require.NoError(t, err)
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

func TestChooseSomeEqualWeights(t *testing.T) {
	r, err := bandrng.NewRng([]byte("THIS_IS_A_RANDOM_SEED_LONG_ENOUGH_FOR_ENTROPY"), []byte("1"), []byte("TEST"))
	require.NoError(t, err)
	length := 10
	weights := make([]uint64, length)
	for idx := 0; idx < length; idx++ {
		weights[idx] = 1
	}
	require.Equal(t, []int{6, 0, 5, 9, 4}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{6, 0, 9, 8, 2}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{5, 6, 1, 8, 3}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{8, 3, 9, 1, 7}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{6, 1, 5, 3, 7}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{4, 9, 3, 1, 5}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{5, 1, 7, 2, 8}, bandrng.ChooseSome(r, weights, 5))
}

func TestChooseSomeSkewedWeights(t *testing.T) {
	r, err := bandrng.NewRng([]byte("THIS_IS_A_RANDOM_SEED_LONG_ENOUGH_FOR_ENTROPY"), []byte("1"), []byte("TEST"))
	require.NoError(t, err)
	length := 10
	weights := make([]uint64, length)
	for idx := 0; idx < length; idx++ {
		weights[idx] = uint64(1 + idx*10)
	}
	require.Equal(t, []int{8, 7, 2, 9, 6}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{6, 8, 7, 4, 3}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{7, 4, 8, 6, 5}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{8, 5, 4, 9, 1}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{6, 1, 7, 9, 8}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{9, 7, 5, 4, 2}, bandrng.ChooseSome(r, weights, 5))
	require.Equal(t, []int{3, 9, 5, 8, 7}, bandrng.ChooseSome(r, weights, 5))
}

func TestChooseSomeMaxWeight(t *testing.T) {
	r, err := bandrng.NewRng([]byte("THIS_IS_A_RANDOM_SEED_LONG_ENOUGH_FOR_ENTROPY"), []byte("1"), []byte("TEST"))
	require.NoError(t, err)
	length := 10
	weights := make([]uint64, length)
	for idx := 0; idx < length; idx++ {
		weights[idx] = uint64(1 + idx*10)
	}
	require.Equal(t, []int{8, 7, 2, 9, 6}, bandrng.ChooseSomeMaxWeight(r, weights, 5, 3))
	require.Equal(t, []int{6, 1, 7, 9, 8}, bandrng.ChooseSomeMaxWeight(r, weights, 5, 3))
	require.Equal(t, []int{9, 3, 8, 7, 6}, bandrng.ChooseSomeMaxWeight(r, weights, 5, 3))
	require.Equal(t, []int{9, 8, 4, 6, 5}, bandrng.ChooseSomeMaxWeight(r, weights, 5, 3))
	require.Equal(t, []int{8, 9, 6, 7, 5}, bandrng.ChooseSomeMaxWeight(r, weights, 5, 3))
	require.Equal(t, []int{3, 9, 7, 6, 8}, bandrng.ChooseSomeMaxWeight(r, weights, 5, 3))
	require.Equal(t, []int{4, 6, 9, 5, 7}, bandrng.ChooseSomeMaxWeight(r, weights, 5, 3))
}
