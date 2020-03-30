package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b *BandDB) AddOrCreateAccount(address sdk.AccAddress, balance sdk.Coins) error {
	return b.tx.Where(Account{Address: address.String()}).
		Assign(Account{Balance: balance.String()}).
		FirstOrCreate(&Account{}).Error
}
