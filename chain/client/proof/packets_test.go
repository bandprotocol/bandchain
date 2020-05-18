package proof

import (
	"encoding/hex"
	"testing"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func TestCalculateResultHash(t *testing.T) {
	// RawByte is d9c589270a046265656210011a0f03000000425443640000000000000020012801
	reqPacket := oracle.OracleRequestPacketData{
		ClientID:       "beeb",
		OracleScriptID: 1,
		Calldata:       mustDecodeString("030000004254436400000000000000"),
		AskCount:       1,
		MinCount:       1,
	}
	// RawByte is 79b5957c0a04626565621003180120acc2f9f50528aec2f9f50530013a084bb10e0000000000
	resPacket := oracle.OracleResponsePacketData{
		ClientID:      "beeb",
		RequestID:     3,
		AnsCount:      1,
		RequestTime:   1589535020,
		ResolveTime:   1589535022,
		ResolveStatus: oracle.ResolveStatus(1),
		Result:        mustDecodeString("4bb10e0000000000"),
	}
	expectedResultHash := hexToBytes("dbbbf5596a975c50c601bdd6ae26a5007e8483344afd7d2ae41e37891cb81b86")

	require.Equal(t, expectedResultHash, tmhash.Sum(
		append(
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(reqPacket)),
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(resPacket))...,
		)))
}

func TestEmptyClientID(t *testing.T) {
	var err error
	calldata, err := hex.DecodeString("030000004254436400000000000000")
	if err != nil {
		panic(err)
	}
	result, err := hex.DecodeString("0aae0e0000000000")
	if err != nil {
		panic(err)
	}
	// RawByte is d9c5892710011a0f03000000425443640000000000000020012801
	reqPacket := oracle.OracleRequestPacketData{
		ClientID:       "",
		OracleScriptID: 1,
		Calldata:       calldata,
		AskCount:       1,
		MinCount:       1,
	}
	// RawByte is 79b5957c1004180120f3caf9f50528f7caf9f50530013a080aae0e0000000000
	resPacket := oracle.OracleResponsePacketData{
		ClientID:      "",
		RequestID:     4,
		AnsCount:      1,
		RequestTime:   1589536115,
		ResolveTime:   1589536119,
		ResolveStatus: oracle.ResolveStatus(1),
		Result:        result,
	}
	expectedResultHash := hexToBytes("37ca0d67535481e7575785dc828b730279f5620ba373c667c22a512cc48ae6dc")

	require.Equal(t, expectedResultHash, tmhash.Sum(
		append(
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(reqPacket)),
			tmhash.Sum(oracle.ModuleCdc.MustMarshalBinaryBare(resPacket))...,
		)))
}
