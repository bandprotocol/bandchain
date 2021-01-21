package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/GeoDB-Limited/odincore/chain/pkg/obi"
)

func mustDecodeString(hexstr string) []byte {
	b, err := hex.DecodeString(hexstr)
	if err != nil {
		panic(err)
	}
	return b
}

func TestGetBytesRequestPacket(t *testing.T) {
	req := OracleRequestPacketData{
		ClientID:       "test",
		OracleScriptID: 1,
		Calldata:       mustDecodeString("030000004254436400000000000000"),
		AskCount:       1,
		MinCount:       1,
	}
	require.Equal(t, []byte(`{"type":"oracle/OracleRequestPacketData","value":{"ask_count":"1","calldata":"AwAAAEJUQ2QAAAAAAAAA","client_id":"test","min_count":"1","oracle_script_id":"1"}}`), req.GetBytes())
}

func TestGetBytesResponsePacket(t *testing.T) {
	res := OracleResponsePacketData{
		ClientID:      "test",
		RequestID:     1,
		AnsCount:      1,
		RequestTime:   1589535020,
		ResolveTime:   1589535022,
		ResolveStatus: ResolveStatus(1),
		Result:        mustDecodeString("4bb10e0000000000"),
	}
	require.Equal(t, []byte(`{"type":"oracle/OracleResponsePacketData","value":{"ans_count":"1","client_id":"test","request_id":"1","request_time":"1589535020","resolve_status":1,"resolve_time":"1589535022","result":"S7EOAAAAAAA="}}`), res.GetBytes())
}

func TestOBIEncodeResult(t *testing.T) {
	req := OracleRequestPacketData{
		ClientID:       "beeb",
		OracleScriptID: 1,
		Calldata:       mustDecodeString("0000000342544300000000000003e8"),
		AskCount:       1,
		MinCount:       1,
	}

	res := OracleResponsePacketData{
		ClientID:      "beeb",
		RequestID:     2,
		AnsCount:      1,
		RequestTime:   1591622616,
		ResolveTime:   1591622618,
		ResolveStatus: ResolveStatus(1),
		Result:        mustDecodeString("00000000009443ee"),
	}
	expectedEncodedResult := mustDecodeString("000000046265656200000000000000010000000f0000000342544300000000000003e800000000000000010000000000000001000000046265656200000000000000020000000000000001000000005ede3bd8000000005ede3bda000000010000000800000000009443ee")
	require.Equal(t, expectedEncodedResult, obi.MustEncode(req, res))
}

func TestOBIEncodeResultOfEmptyClientID(t *testing.T) {
	req := OracleRequestPacketData{
		ClientID:       "",
		OracleScriptID: 1,
		Calldata:       mustDecodeString("0000000342544300000000000003e8"),
		AskCount:       1,
		MinCount:       1,
	}

	res := OracleResponsePacketData{
		ClientID:      "",
		RequestID:     1,
		AnsCount:      1,
		RequestTime:   1591622426,
		ResolveTime:   1591622429,
		ResolveStatus: ResolveStatus(1),
		Result:        mustDecodeString("0000000000944387"),
	}
	expectedEncodedResult := mustDecodeString("0000000000000000000000010000000f0000000342544300000000000003e8000000000000000100000000000000010000000000000000000000010000000000000001000000005ede3b1a000000005ede3b1d00000001000000080000000000944387")
	require.Equal(t, expectedEncodedResult, obi.MustEncode(req, res))
}
