package keeper

import (
	"github.com/bandprotocol/d3n/chain/x/zoracle/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) CheckReporter(ctx sdk.Context, validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.ReporterStoreKey(validatorAddress, reporterAddress))
}

func (k Keeper) AddReporter(
	ctx sdk.Context, validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress,
) sdk.Error {
	if k.CheckReporter(ctx, validatorAddress, reporterAddress) {
		return types.ErrInvalidState(
			"AddReporter:  (%s) is already a reporter of (%s).",
			reporterAddress.String(),
			validatorAddress.String(),
		)
	}

	store := ctx.KVStore(k.storeKey)
	store.Set(types.ReporterStoreKey(validatorAddress, reporterAddress), []byte{1})

	return nil
}

func (k Keeper) RemoveReporter(
	ctx sdk.Context, validatorAddress sdk.ValAddress, reporterAddress sdk.AccAddress,
) sdk.Error {
	if !k.CheckReporter(ctx, validatorAddress, reporterAddress) {
		return types.ErrItemNotFound(
			"RemoveReporter: Item notfound for validator address (%s) reporter address (%s).",
			validatorAddress.String(),
			reporterAddress.String(),
		)
	}

	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ReporterStoreKey(validatorAddress, reporterAddress))

	return nil
}
