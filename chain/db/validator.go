package db

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/tendermint/tendermint/libs/common"
)

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
	validatorStruct := Validator{OperatorAddress: validator.String()}
	err := b.tx.First(&validatorStruct).Error
	return validatorStruct, err == nil
}

func (b *BandDB) handleMsgEditValidator(msg staking.MsgEditValidator) error {
	validator, isFound := b.GetValidator(msg.ValidatorAddress)
	if !isFound {
		return fmt.Errorf(fmt.Sprintf("validator %s has not exist.", msg.ValidatorAddress.String()))
	}

	if msg.Description.Moniker != staking.DoNotModifyDesc {
		validator.Moniker = msg.Description.Moniker
	}
	if msg.Description.Identity != staking.DoNotModifyDesc {
		validator.Identity = msg.Description.Identity
	}
	if msg.Description.Website != staking.DoNotModifyDesc {
		validator.Website = msg.Description.Website
	}
	if msg.Description.Details != staking.DoNotModifyDesc {
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
	return b.tx.Create(&Validator{
		OperatorAddress:     msg.ValidatorAddress.String(),
		ConsensusAddress:    msg.PubKey.Address().String(),
		Moniker:             msg.Description.Moniker,
		Identity:            msg.Description.Identity,
		Website:             msg.Description.Website,
		Details:             msg.Description.Details,
		CommissionRate:      msg.Commission.Rate.String(),
		CommissionMaxRate:   msg.Commission.MaxRate.String(),
		CommissionMaxChange: msg.Commission.MaxChangeRate.String(),
		MinSelfDelegation:   msg.MinSelfDelegation.ToDec().String(),
		SelfDelegation:      msg.Value.Amount.ToDec().String(),
	}).Error
}
