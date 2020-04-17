package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

// AnyError returns the first error found in the given error list, or nil if none exists.
func AnyError(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

// EnsureMaxValue checks whether the given value exceeds the max limit as specified by param.
func (k Keeper) EnsureMaxValue(ctx sdk.Context, param []byte, val uint64) error {
	maxVal := k.GetParam(ctx, param)
	if val > maxVal {
		return sdkerrors.Wrapf(types.ErrBadDataLength,
			"%s violation; got: %d, max: %d", string(param), val, maxVal,
		)
	}
	return nil
}
