package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// ContainsVal returns whether the given slice of validators contains the target validator.
func ContainsVal(vals []sdk.ValAddress, target sdk.ValAddress) bool {
	for _, val := range vals {
		if val.Equals(target) {
			return true
		}
	}
	return false
}

// ContainsEID returns whether the given slice of external ids contains the target id.
func ContainsEID(ids []types.ExternalID, target types.ExternalID) bool {
	for _, id := range ids {
		if id == target {
			return true
		}
	}
	return false
}
