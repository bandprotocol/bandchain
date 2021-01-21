package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
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

// ContainsEID returns whether the given slice of raw requests contains the target id.
func ContainsEID(rawRequests []types.RawRequest, target types.ExternalID) bool {
	for _, rawRequest := range rawRequests {
		if rawRequest.ExternalID == target {
			return true
		}
	}
	return false
}
