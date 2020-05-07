package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func (k Keeper) AfterValidatorBonded(ctx sdk.Context, _ sdk.ConsAddress, address sdk.ValAddress) {
	// Create a new report info if not existed.
	fmt.Println("DDDDDD")
	found := k.HasValidatorReportInfo(ctx, address)
	if !found {
		reportInfo := types.NewValidatorReportInfo(
			address,
			false,
			0,
			0,
		)
		k.SetValidatorReportInfo(ctx, address, reportInfo)
	}
}

//_________________________________________________________________________________________

// Hooks wrapper struct for slashing keeper
type Hooks struct {
	k Keeper
}

var _ types.StakingHooks = Hooks{}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Implements sdk.ValidatorHooks
func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) {
	h.k.AfterValidatorBonded(ctx, consAddr, valAddr)
}

func (h Hooks) AfterValidatorRemoved(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)         {}
func (h Hooks) AfterValidatorCreated(_ sdk.Context, _ sdk.ValAddress)                            {}
func (h Hooks) AfterValidatorBeginUnbonding(_ sdk.Context, _ sdk.ConsAddress, _ sdk.ValAddress)  {}
func (h Hooks) BeforeValidatorModified(_ sdk.Context, _ sdk.ValAddress)                          {}
func (h Hooks) BeforeDelegationCreated(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeDelegationSharesModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress) {}
func (h Hooks) BeforeDelegationRemoved(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) AfterDelegationModified(_ sdk.Context, _ sdk.AccAddress, _ sdk.ValAddress)        {}
func (h Hooks) BeforeValidatorSlashed(_ sdk.Context, _ sdk.ValAddress, _ sdk.Dec)                {}
