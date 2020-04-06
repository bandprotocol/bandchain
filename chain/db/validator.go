package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

func (b *BandDB) AddValidatorUpTime(
	rawConsensusAddress tmbytes.HexBytes,
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

func (b *BandDB) CreateValidator(
	validatorAddress sdk.ValAddress,
	pubKey crypto.PubKey,
	moniker string,
	identity string,
	website string,
	details string,
	commissionRate sdk.Dec,
	commissionMaxRate sdk.Dec,
	commissionMaxChangeRate sdk.Dec,
	minSelfDelegation sdk.Int,
	value sdk.Coin,
) error {
	return b.tx.Create(&Validator{
		OperatorAddress:     validatorAddress.String(),
		ConsensusAddress:    pubKey.Address().String(),
		Moniker:             moniker,
		Identity:            identity,
		Website:             website,
		Details:             details,
		CommissionRate:      commissionRate.String(),
		CommissionMaxRate:   commissionMaxRate.String(),
		CommissionMaxChange: commissionMaxChangeRate.String(),
		MinSelfDelegation:   minSelfDelegation.String(),
		Tokens:              value.Amount.String(),
		DelegatorShares:     value.Amount.String(),
	}).Error
}

func (b *BandDB) UpdateValidator(validatorAddress sdk.ValAddress, details *Validator) error {
	return b.tx.Model(&Validator{}).
		Where(Validator{OperatorAddress: validatorAddress.String()}).
		Update(details).Error
}

func (b *BandDB) RemoveValidator(validatorAddress sdk.ValAddress) error {
	return b.tx.Delete(Validator{
		OperatorAddress: validatorAddress.String(),
	}).Error
}

func (b *BandDB) GetValidator(validator sdk.ValAddress) (Validator, bool) {
	validatorStruct := Validator{OperatorAddress: validator.String()}
	err := b.tx.First(&validatorStruct).Error
	return validatorStruct, err == nil
}
