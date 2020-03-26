package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
