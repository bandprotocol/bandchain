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
	// TODO
	return nil
}

func (b *BandDB) handleMsgDelegate(msg staking.MsgDelegate) error {
	// TODO
	return nil
}
