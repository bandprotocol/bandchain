package keeper

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IsReporter returns true iff the address is an authorized reporter for the given validator.
func (k Keeper) IsReporter(ctx sdk.Context, val sdk.ValAddress, addr sdk.AccAddress) bool {
	if val.Equals(sdk.ValAddress(addr)) { // A validator is always a reporter of himself.
		return true
	}
	return ctx.KVStore(k.storeKey).Has(types.ReporterStoreKey(val, addr))
}

// AddReporter adds the reporter address to the list of reporters of the given validator.
func (k Keeper) AddReporter(ctx sdk.Context, val sdk.ValAddress, addr sdk.AccAddress) error {
	if k.IsReporter(ctx, val, addr) {
		return sdkerrors.Wrapf(
			types.ErrReporterAlreadyExists, "val: %s, addr: %s", val.String(), addr.String())
	}
	ctx.KVStore(k.storeKey).Set(types.ReporterStoreKey(val, addr), []byte{1})
	return nil
}

// RemoveReporter removes the reporter address from the list of reporters of the given validator.
func (k Keeper) RemoveReporter(ctx sdk.Context, val sdk.ValAddress, addr sdk.AccAddress) error {
	if !k.IsReporter(ctx, val, addr) {
		return sdkerrors.Wrapf(
			types.ErrReporterNotFound, "val: %s, addr: %s", val.String(), addr.String())
	}
	ctx.KVStore(k.storeKey).Delete(types.ReporterStoreKey(val, addr))
	return nil
}

// GetAllReportersOfValidator returns the reporter list of the given validator.
func (k Keeper) GetAllReportersOfValidator(ctx sdk.Context, val sdk.ValAddress) (reporters []sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	storeKeyPrefix := append(types.ReporterStoreKeyPrefix, val.Bytes()...)
	iterator := sdk.KVStorePrefixIterator(store, storeKeyPrefix)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		key := iterator.Key()
		reporterBytes := key[21:]
		reporterAddress := sdk.AccAddress(reporterBytes)

		reporters = append(reporters, reporterAddress)
	}
	return reporters
}
