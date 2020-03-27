package db

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
)

func (b *BandDB) AddValidator(
	operatorAddress sdk.ValAddress,
	consensusAddress crypto.PubKey,
	moniker string,
	identity string,
	website string,
	details string,
	commissionRate string,
	commissionMaxRate string,
	commissionMaxChange string,
	minSelfDelegation string,
	selfDelegation string,
) error {
	return b.tx.Create(&Validator{
		OperatorAddress:     operatorAddress.String(),
		ConsensusAddress:    consensusAddress.Address().String(),
		Moniker:             moniker,
		Identity:            identity,
		Website:             website,
		Details:             details,
		CommissionRate:      commissionRate,
		CommissionMaxRate:   commissionMaxRate,
		CommissionMaxChange: commissionMaxChange,
		MinSelfDelegation:   minSelfDelegation,
		SelfDelegation:      selfDelegation,
	}).Error
}

func (b *BandDB) AddValidatorUpTime(
	rawConsensusAddress common.HexBytes,
	height int64,
	voted bool,
) error {
	consensusAddress := rawConsensusAddress.String()
	err := b.tx.Create(&ValidatorVote{
		ConsensusAddress: consensusAddress,
		BlockHeight:      height,
		Voted:            voted,
	}).Error

	if err != nil {
		return err
	}

	var validator Validator
	err = b.tx.Where(Validator{ConsensusAddress: consensusAddress}).First(&validator).Error
	if err != nil {
		return err
	}

	validator.ElectedCount++
	if voted {
		validator.VotedCount++
	} else {
		validator.MissedCount++
	}

	b.tx.Save(&validator)
	return nil
}

func (b *BandDB) ClearOldVotes(currentHeight int64) error {
	uptimeLookBackDuration, err := b.GetUptimeLookBackDuration()
	if err != nil {
		return err
	}

	if currentHeight > uptimeLookBackDuration {
		var votes []ValidatorVote
		err := b.tx.Find(
			&votes,
			"block_height <= ?",
			currentHeight-uptimeLookBackDuration,
		).Error

		if err != nil {
			return err
		}
		for _, vote := range votes {
			var validator Validator
			err = b.tx.Where(Validator{ConsensusAddress: vote.ConsensusAddress}).First(&validator).Error
			if err == nil {
				validator.ElectedCount--
				if vote.Voted {
					validator.VotedCount--
				} else {
					validator.MissedCount--
				}
				b.tx.Save(&validator)
			}

		}
		return b.tx.Delete(
			ValidatorVote{},
			"block_height <= ?",
			currentHeight-uptimeLookBackDuration,
		).Error
	}
	return nil
}

func (b *BandDB) GetValidator(validator sdk.ValAddress) (Validator, bool) {
	validatorSturct := Validator{OperatorAddress: validator.String()}
	err := b.tx.First(&validatorSturct).Error
	if err == nil {
		// Found a validator
		return validatorSturct, true
	}
	// Not found any validator
	return validatorSturct, false
}

func (b *BandDB) handleMsgEditValidator(msg staking.MsgEditValidator) error {
	validator, err := b.GetValidator(msg.ValidatorAddress)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("validator %s has not exist.", msg.ValidatorAddress.String()))
	}

	if msg.Description.Moniker != "[do-not-modify]" {
		validator.Moniker = msg.Description.Moniker
	}
	if msg.Description.Identity != "[do-not-modify]" {
		validator.Identity = msg.Description.Identity
	}
	if msg.Description.Website != "[do-not-modify]" {
		validator.Website = msg.Description.Website
	}
	if msg.Description.Details != "[do-not-modify]" {
		validator.Details = msg.Description.Details
	}
	if msg.CommissionRate != nil {
		validator.CommissionRate = msg.CommissionRate.String()
	}
	if msg.MinSelfDelegation != nil {
		validator.MinSelfDelegation = msg.MinSelfDelegation.ToDec().String()
	}

	return b.tx.Save(&validator).Error
}

func (b *BandDB) handleMsgCreateValidator(msg staking.MsgCreateValidator) error {
	_, err := b.GetValidator(msg.ValidatorAddress)
	if err == nil {
		return fmt.Errorf(fmt.Sprintf("validator %s has already exist.", msg.ValidatorAddress.String()))
	}

	return b.AddValidator(
		msg.ValidatorAddress,
		msg.PubKey,
		msg.Description.Moniker,
		msg.Description.Identity,
		msg.Description.Website,
		msg.Description.Details,
		msg.Commission.Rate.String(),
		msg.Commission.MaxRate.String(),
		msg.Commission.MaxChangeRate.String(),
		msg.MinSelfDelegation.ToDec().String(),
		msg.Value.Amount.ToDec().String(),
	)
}
