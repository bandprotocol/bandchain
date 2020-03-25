package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b *BandDB) AddDelegation(
	delegatorAddress sdk.AccAddress,
	operatorAddress sdk.ValAddress,
	Amount sdk.Coin,
) error {
	return b.tx.Create(&Delegation{
		DelegatorAddress: delegatorAddress.String(),
		OperatorAddress:  operatorAddress.String(),
		Amount:           Amount.String(),
	}).Error
}
