package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b *BandDB) SetAccountBalance(address sdk.AccAddress, balance sdk.Coins) error {
	return b.tx.Where(Account{Address: address.String()}).
		Assign(Account{Balance: balance.String()}).
		FirstOrCreate(&Account{}).Error
}

func (b *BandDB) DecreaseAccountBalance(address sdk.AccAddress, balance sdk.Coins) error {
	account := Account{Address: address.String()}
	err := b.tx.First(&account).Error
	if err != nil {
		return err
	}
	currentBalance, err := sdk.ParseCoins(account.Balance)
	if err != nil {
		return err
	}
	account.Balance = currentBalance.Sub(balance).String()
	return b.tx.Save(&account).Error
}
