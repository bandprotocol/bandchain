package types

import (
	"testing"

	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
	"github.com/stretchr/testify/require"
)

func TestLatestResponseStoreKey(t *testing.T) {
	request := otypes.NewOracleRequestPacketData("alice", 1, []byte("calldata"), 1, 1)
	expected := append([]byte{0x03}, request.GetBytes()...)
	require.Equal(t, expected, LastestResponseStoreKey(request))
}
