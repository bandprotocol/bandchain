package db

import (
	"encoding/json"
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
	messages json.RawMessage,
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
	messages []map[string]interface{},
) error {
	rawJson, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	transaction := createTransaction(
		txHash,
		timestamp,
		gasUsed,
		gasLimit,
		gasFee,
		sender,
		success,
		blockHeight,
		rawJson,
	)
	err = b.tx.Create(&transaction).Error
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
