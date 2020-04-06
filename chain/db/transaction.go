package db

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createTransaction(
	id int64,
	txHash []byte,
	timestamp time.Time,
	gasUsed int64,
	gasLimit uint64,
	gasFee sdk.Coins,
	sender sdk.AccAddress,
	success bool,
	blockHeight int64,
) Transaction {
	rawJson, _ := json.Marshal(nil)
	return Transaction{
		ID:          id,
		TxHash:      txHash,
		Timestamp:   timestamp.UnixNano() / int64(time.Millisecond),
		GasUsed:     gasUsed,
		GasLimit:    gasLimit,
		GasFee:      gasFee.String(),
		Sender:      sender.String(),
		Success:     success,
		BlockHeight: blockHeight,
		Messages:    rawJson,
	}
}

func (b *BandDB) AddTransaction(
	id int64,
	txHash []byte,
	timestamp time.Time,
	gasUsed int64,
	gasLimit uint64,
	gasFee sdk.Coins,
	sender sdk.AccAddress,
	success bool,
	blockHeight int64,
) error {
	transaction := createTransaction(
		id,
		txHash,
		timestamp,
		gasUsed,
		gasLimit,
		gasFee,
		sender,
		success,
		blockHeight,
	)

	err := b.tx.Create(&transaction).Error
	return err
}

func (b *BandDB) UpdateTransaction(
	txHash []byte,
	messages []map[string]interface{},
) error {

	var transaction Transaction
	err := b.tx.First(&transaction, txHash).Error

	if err != nil {
		return err
	}
	rawJson, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	transaction.Messages = rawJson
	err = b.tx.Save(&transaction).Error

	return err
}
