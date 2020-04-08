package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CheckReporter returns true iff the given reporter is authorized to report data on behalf of
// the given validator.
func (k Keeper) CheckReporter(
	ctx sdk.Context, validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress,
) bool {
	if validatorAddress.Equals(sdk.ValAddress(reporterAddress)) {
		return true
	}
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ReporterStoreKey(validatorAddress, reporterAddress))
}

// AddReporter adds the given reporter to the list of reporters of the given validator.
func (k Keeper) AddReporter(
	ctx sdk.Context, validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress,
) error {
	if k.CheckReporter(ctx, validatorAddress, reporterAddress) {
		return sdkerrors.Wrapf(types.ErrItemDuplication,
			"AddReporter: %s is already a reporter of %s.",
			reporterAddress.String(), validatorAddress.String(),
		)
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ReporterStoreKey(validatorAddress, reporterAddress), []byte{1})
	return nil
}

// AddReporter removes the given reporter from the list of reporters of the given validator.
func (k Keeper) RemoveReporter(
	ctx sdk.Context, validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress,
) error {
	if !k.CheckReporter(ctx, validatorAddress, reporterAddress) {
		return sdkerrors.Wrapf(types.ErrItemNotFound,
			"RemoveReporter: %s is not a reporter of %s.",
			reporterAddress.String(), validatorAddress.String(),
		)
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ReporterStoreKey(validatorAddress, reporterAddress))
	return nil
}
