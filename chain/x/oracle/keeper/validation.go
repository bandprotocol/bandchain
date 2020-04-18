package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

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

// AnyError returns the first error found in the given error list, or nil if none exists.
func AnyError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// EnsureMax checks whether the given uint64 value exceeds the max limit as specified by param.
func (k Keeper) EnsureMax(ctx sdk.Context, param []byte, val uint64) error {
	maxVal := k.GetParam(ctx, param)
	if val > maxVal {
		return sdkerrors.Wrapf(types.ErrBadDataLength,
			"%s violation; got: %d, max: %d", string(param), val, maxVal,
		)
	}
	return nil
}

// EnsureLength checks whether the given int value exceeds the max limit as specified by param.
func (k Keeper) EnsureLength(ctx sdk.Context, param []byte, len int) error {
	return k.EnsureMax(ctx, param, uint64(len))
}
