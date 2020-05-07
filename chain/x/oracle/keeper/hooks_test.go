package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func TestAfterValidatorBonded(t *testing.T) {
	_, ctx, k := createTestInput()

	_, err := k.GetValidatorReportInfo(ctx, Alice.ValAddress)
	require.Error(t, err)

	info, err := k.GetValidatorReportInfo(ctx, Validator1.ValAddress)
	require.Nil(t, err)
	require.Equal(t, types.NewValidatorReportInfo(Validator1.ValAddress, false, 0, 0), info)
}
