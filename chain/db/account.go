package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (b *BandDB) SetAccountBalance(
	address sdk.AccAddress, balance sdk.Coins, blockHeight int64,
) error {
	return b.tx.Where(Account{Address: address.String()}).
		Assign(Account{Balance: balance.String(), UpdatedHeight: blockHeight}).
		FirstOrCreate(&Account{}).Error
}

func (b *BandDB) DecreaseAccountBalance(
	address sdk.AccAddress, balance sdk.Coins, blockHeight int64,
) error {
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
	account.UpdatedHeight = blockHeight
	return b.tx.Save(&account).Error
}
