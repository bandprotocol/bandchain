package keeper_test
import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MockRequest()types.Request{
	oracleScriptID := types.OracleScriptID(1)
	calldata := []byte("CALLDATA")
	requestedValidators := []sdk.ValAddress{Validator1.ValAddress,Validator2.ValAddress,Validator3.ValAddress}
	minCount := int64(2)
	requestHeight := int64(20)
	requestTime := int64(1589535020)
	clientID := "beeb"
	ibcInfo := types.NewIBCInfo("source_port","source_channel")
	rawRequestIDs := []types.ExternalID{1,2,3}
	return types.NewRequest(
		oracleScriptID,calldata,requestedValidators,minCount,
		requestHeight,requestTime,clientID,&ibcInfo,rawRequestIDs,
	)
}


func TestPrepareRequestSuccess(t *testing.T) {
	_, ctx, k := createTestInput()
	
	// os, clear:= loadOracleScriptFromWasm()
	fmt.Println(Validator1.ValAddress,Validator2.ValAddress,Validator3.ValAddress)
	
	oracleScriptID := types.OracleScriptID(1)
	calldata := []byte("CALLDATA")
	askCount := int64(1)
	minCount := int64(2)
	clientID := "beeb"
	sender := Alice.Address
	ibcInfo := types.NewIBCInfo("source_port","source_channel")
	
	m := types.NewMsgRequestData(oracleScriptID,calldata,askCount,minCount,clientID,sender)
	err := k.PrepareRequest(ctx,&m,&ibcInfo)
	require.NoError(t,err)
}