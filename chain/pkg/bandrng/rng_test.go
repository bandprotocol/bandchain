package bandrng_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/pkg/bandrng"
)

func TestRngRandom(t *testing.T) {
	r := bandrng.NewRng("SEED")
	require.Equal(t, r.NextUint64(), uint64(15735084640102210068))
	require.Equal(t, r.NextUint64(), uint64(3485776390957061973))
	require.Equal(t, r.NextUint64(), uint64(17609118114147816341))
	require.Equal(t, r.NextUint64(), uint64(15960811988050104523))
	require.Equal(t, r.NextUint64(), uint64(11919533627209787235))
}
