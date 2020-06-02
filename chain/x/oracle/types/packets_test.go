package types

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
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

func TestCalculateResultHash(t *testing.T) {
	req := OracleRequestPacketData{
		ClientID:       "beeb",
		OracleScriptID: 1,
		Calldata:       mustDecodeString("030000004254436400000000000000"),
		AskCount:       1,
		MinCount:       1,
	}

	res := OracleResponsePacketData{
		ClientID:      "beeb",
		RequestID:     1,
		AnsCount:      1,
		RequestTime:   1589535020,
		ResolveTime:   1589535022,
		ResolveStatus: ResolveStatus(1),
		Result:        mustDecodeString("4bb10e0000000000"),
	}
	expectedResultHash := mustDecodeString("29bcc52d59b39c61a9616365dcd39dbff8d1aebc88a6a7e2b53dff67841dbc06")
	require.Equal(t, expectedResultHash, CalculateResultHash(req, res))
}

func TestCalculateResultHashOfEmptyClientID(t *testing.T) {
	req := OracleRequestPacketData{
		ClientID:       "",
		OracleScriptID: 1,
		Calldata:       mustDecodeString("030000004254436400000000000000"),
		AskCount:       1,
		MinCount:       1,
	}

	res := OracleResponsePacketData{
		ClientID:      "",
		RequestID:     1,
		AnsCount:      1,
		RequestTime:   1590490752,
		ResolveTime:   1590490756,
		ResolveStatus: ResolveStatus(1),
		Result:        mustDecodeString("568c0d0000000000"),
	}
	expectedResultHash := mustDecodeString("a506eb6a23931d1130bfced8b10ec41674a6eaefb888063847e3605bffbdd5ba")
	require.Equal(t, expectedResultHash, CalculateResultHash(req, res))
}
