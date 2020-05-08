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

	b.tx.Save(&validator)
	return nil
}

func (b *BandDB) CreateValidator(
	validatorAddress sdk.ValAddress,
	pubkey crypto.PubKey,
	moniker string,
	identity string,
	website string,
	details string,
	commissionRate sdk.Dec,
	commissionMaxRate sdk.Dec,
	commissionMaxChangeRate sdk.Dec,
	minSelfDelegation sdk.Int,
	value sdk.Coin,
	bondedHeight int64,
) error {
	tokens := value.Amount.Uint64()
	jailed := false
	return b.tx.Create(&Validator{
		OperatorAddress:     validatorAddress.String(),
		ConsensusAddress:    pubkey.Address().String(),
		ConsensusPubkey:     sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubkey),
		Moniker:             moniker,
		Identity:            identity,
		Website:             website,
		Details:             details,
		CommissionRate:      commissionRate.String(),
		CommissionMaxRate:   commissionMaxRate.String(),
		CommissionMaxChange: commissionMaxChangeRate.String(),
		MinSelfDelegation:   minSelfDelegation.String(),
		Tokens:              &tokens,
		DelegatorShares:     value.Amount.String(),
		BondedHeight:        bondedHeight,
		CurrentReward:       "0",
		CurrentRatio:        "0",
		Jailed:              &jailed,
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
