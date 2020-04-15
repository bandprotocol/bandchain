package db

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b *BandDB) SetDelegation(
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
	shares string,
) error {
	return b.tx.Where(Delegation{
		DelegatorAddress: delegatorAddress.String(),
		ValidatorAddress: validatorAddress.String(),
	}).
		Assign(Delegation{Shares: shares}).
		FirstOrCreate(&Delegation{}).Error
}

func (b *BandDB) RemoveDelegation(
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
) error {
	return b.tx.Delete(Delegation{
		DelegatorAddress: delegatorAddress.String(),
		ValidatorAddress: validatorAddress.String(),
	}).Error
}

func (b *BandDB) delegate(
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
) error {
	validator, found := b.StakingKeeper.GetValidator(b.ctx, validatorAddress)
	if !found {
		return errors.New("Not found validator")
	}
	err := b.UpdateValidator(
		validatorAddress,
		&Validator{
			Tokens:          validator.Tokens.Uint64(),
			DelegatorShares: validator.DelegatorShares.String(),
		},
	)
	if err != nil {
		return err
	}

	delegation, found := b.StakingKeeper.GetDelegation(
		b.ctx, delegatorAddress, validatorAddress,
	)
	if !found {
		return errors.New("Not found delegation")
	}
	return b.SetDelegation(delegatorAddress, validatorAddress, delegation.Shares.String())
}

func (b *BandDB) undelegate(
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
) error {
	delegation, found := b.StakingKeeper.GetDelegation(
		b.ctx, delegatorAddress, validatorAddress,
	)
	if found {
		err := b.SetDelegation(
			delegatorAddress,
			validatorAddress,
			delegation.Shares.String(),
		)
		if err != nil {
			return err
		}
	} else {
		err := b.RemoveDelegation(delegatorAddress, validatorAddress)
		if err != nil {
			return err
		}
	}

	validator, found := b.StakingKeeper.GetValidator(b.ctx, validatorAddress)
	if found {
		return b.UpdateValidator(
			validatorAddress,
			&Validator{
				Tokens:          validator.Tokens.Uint64(),
				DelegatorShares: validator.DelegatorShares.String(),
				Jailed:          validator.Jailed,
			},
		)
	} else {
		return b.RemoveValidator(validatorAddress)
	}
}
