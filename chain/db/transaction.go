package db

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createTransaction(
	txHash []byte,
	timestamp time.Time,
	gasUsed int64,
	gasLimit uint64,
	gasFee sdk.Coins,
	sender sdk.AccAddress,
	success bool,
	blockHeight int64,
	messages string,
) Transaction {
	return Transaction{
		TxHash:      txHash,
		Timestamp:   timestamp.UnixNano() / int64(time.Millisecond),
		GasUsed:     gasUsed,
		GasLimit:    gasLimit,
		GasFee:      gasFee.String(),
		Sender:      sender.String(),
		Success:     success,
		BlockHeight: blockHeight,
		Messages:    messages,
	}
}

func (b *BandDB) AddTransaction(
	txHash []byte,
	timestamp time.Time,
	gasUsed int64,
	gasLimit uint64,
	gasFee sdk.Coins,
	sender sdk.AccAddress,
	success bool,
	blockHeight int64,
	messages string,
) error {
	transaction := createTransaction(
		txHash,
		timestamp,
		gasUsed,
		gasLimit,
		gasFee,
		sender,
		success,
		blockHeight,
		messages,
	)
	err := b.tx.Create(&transaction).Error
	return err
}

func (b *BandDB) UpdateTransaction(
	txHash []byte,
	messages string,
) error {

	var transaction Transaction
	err := b.tx.First(&transaction, txHash).Error

	if err != nil {
		return err
	}

	transaction.Messages = messages
	err = b.tx.Save(&transaction).Error

	return err
}
