package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

func (b *BandDB) AddDelegation(
	delegatorAddress sdk.AccAddress,
	operatorAddress sdk.ValAddress,
	Coin sdk.Coin,
) error {
	return b.tx.Create(&Delegation{
		DelegatorAddress: delegatorAddress.String(),
		OperatorAddress:  operatorAddress.String(),
		Amount:           Coin.Amount.ToDec().String(),
	}).Error
}

func (b *BandDB) handleMsgDelegate(msg staking.MsgDelegate) error {
	delegation := Delegation{DelegatorAddress: msg.DelegatorAddress.String(), OperatorAddress: msg.ValidatorAddress.String()}
	err := b.tx.First(&delegation).Error
	if err != nil {
		return b.AddDelegation(msg.DelegatorAddress, msg.ValidatorAddress, msg.Amount)
	}
	accumulationAmount, err := sdk.NewDecFromStr(delegation.Amount)
	if err != nil {
		return err
	}

	delegation.Amount = accumulationAmount.Add(msg.Amount.Amount.ToDec()).String()

	return b.tx.Save(&delegation).Error
}
