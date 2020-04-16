package db

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b *BandDB) SetDelegation(
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
	shares string,
	lastRatio sdk.DecCoins,
) error {
	value := "0"
	if !lastRatio.IsZero() {
		value = lastRatio[0].Amount.String()
	}
	return b.tx.Where(Delegation{
		DelegatorAddress: delegatorAddress.String(),
		ValidatorAddress: validatorAddress.String(),
	}).
		Assign(Delegation{Shares: shares, LastRatio: value}).
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
	info := b.DistrKeeper.GetDelegatorStartingInfo(b.ctx, validatorAddress, delegatorAddress)
	latestReward := b.DistrKeeper.GetValidatorHistoricalRewards(b.ctx, validatorAddress, info.PreviousPeriod)
	// CurrentReward must be reset after delegation.
	cumulativeRewardRatio := "0"
	if !latestReward.CumulativeRewardRatio.IsZero() {
		cumulativeRewardRatio = latestReward.CumulativeRewardRatio[0].Amount.String()
	}
	err := b.UpdateValidator(
		validatorAddress,
		&Validator{
			Tokens:          validator.Tokens.BigInt().Uint64(),
			DelegatorShares: validator.DelegatorShares.String(),
			CurrentReward:   "0",
			CurrentRatio:    cumulativeRewardRatio,
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
	return b.SetDelegation(delegatorAddress, validatorAddress, delegation.Shares.String(), latestReward.CumulativeRewardRatio)
}

func (b *BandDB) undelegate(
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
) error {
	delegation, found := b.StakingKeeper.GetDelegation(
		b.ctx, delegatorAddress, validatorAddress,
	)
	if found {
		info := b.DistrKeeper.GetDelegatorStartingInfo(b.ctx, validatorAddress, delegatorAddress)
		latestReward := b.DistrKeeper.GetValidatorHistoricalRewards(b.ctx, validatorAddress, info.PreviousPeriod)
		err := b.SetDelegation(
			delegatorAddress,
			validatorAddress,
			delegation.Shares.String(),
			latestReward.CumulativeRewardRatio,
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
		reward := b.DistrKeeper.GetValidatorCurrentRewards(b.ctx, validatorAddress)
		latestReward := b.DistrKeeper.GetValidatorHistoricalRewards(b.ctx, validatorAddress, reward.Period-1)

		cumulativeRewardRatio := "0"
		if !latestReward.CumulativeRewardRatio.IsZero() {
			cumulativeRewardRatio = latestReward.CumulativeRewardRatio[0].Amount.String()
		}
		return b.UpdateValidator(
			validatorAddress,
			&Validator{
				Tokens:          validator.Tokens.BigInt().Uint64(),
				DelegatorShares: validator.DelegatorShares.String(),
				Jailed:          validator.Jailed,
				CurrentReward:   "0",
				CurrentRatio:    cumulativeRewardRatio,
			},
		)
	} else {
		return b.RemoveValidator(validatorAddress)
	}
}
