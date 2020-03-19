package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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
) sdk.Error {
	if k.CheckReporter(ctx, validatorAddress, reporterAddress) {
		return types.ErrItemDuplication(
			"AddReporter: %s is already a reporter of %s.",
			reporterAddress.String(),
			validatorAddress.String(),
		)
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.ReporterStoreKey(validatorAddress, reporterAddress), []byte{1})

	return nil
}

// AddReporter removes the given reporter from the list of reporters of the given validator.
func (k Keeper) RemoveReporter(
	ctx sdk.Context, validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress,
) sdk.Error {
	if !k.CheckReporter(ctx, validatorAddress, reporterAddress) {
		return types.ErrItemNotFound(
			"RemoveReporter: %s is not a reporter of %s.",
			reporterAddress.String(),
			validatorAddress.String(),
		)
	}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ReporterStoreKey(validatorAddress, reporterAddress))
	return nil
}
