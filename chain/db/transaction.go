package db

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createTransaction(
	txHash []byte,
	timestamp time.Time,
	gasUse uint64,
	gasLimit uint64,
	gasFee string,
	sender sdk.AccAddress,
	success bool,
	blockHeight int64,
) Transaction {
	return Transaction{
		TxHash:      txHash,
		Timestamp:   timestamp,
		GasUse:      gasUse,
		GasLimit:    gasLimit,
		GasFee:      gasFee,
		Sender:      sender.String(),
		Success:     success,
		BlockHeight: blockHeight,
	}
}

func (b *BandDB) AddTransaction(
	txHash []byte,
	timestamp time.Time,
	gasUse uint64,
	gasLimit uint64,
	gasFee string,
	sender sdk.AccAddress,
	success bool,
	blockHeight int64,
) error {
	transaction := createTransaction(
		txHash,
		timestamp,
		gasUse,
		gasLimit,
		gasFee,
		sender,
		success,
		blockHeight,
	)
	err := b.tx.Create(&transaction).Error
	return err
}

// type Transaction struct {
// 	TxHash      []byte `gorm:"primary_key"`
// 	Timestamp   time.Time
// 	GasUse      uint64
// 	GasLimit    uint64
// 	GasFee      sdk.Coins
// 	Sender      sdk.AccAddress
// 	Status      string
// 	BlockHeight int64
// }

// func (b *BandDB) handleMsgEditDataSource(
// 	txHash []byte,
// 	msg zoracle.MsgEditDataSource,
// 	events map[string]string,
// ) error {
// 	dataSource := createDataSource(
// 		int64(msg.DataSourceID), msg.Name, msg.Description,
// 		msg.Owner, msg.Fee, msg.Executable, b.ctx.BlockTime(),
// 	)

// 	err := b.tx.Save(&dataSource).Error
// 	if err != nil {
// 		return err
// 	}

// 	return b.tx.Create(&DataSourceRevision{
// 		DataSourceID: int64(msg.DataSourceID),
// 		Name:         msg.Name,
// 		Timestamp:    b.ctx.BlockTime(),
// 		BlockHeight:  b.ctx.BlockHeight(),
// 		TxHash:       txHash,
// 	}).Error
// }
