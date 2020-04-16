package db

// import (
// 	"errors"
// 	"strconv"
// 	"time"

// 	"github.com/bandprotocol/bandchain/chain/x/oracle"
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// )

// func createDataSource(
// 	id int64,
// 	name, description string,
// 	owner sdk.AccAddress,
// 	fee sdk.Coins,
// 	executable []byte,
// 	blockTime time.Time,
// ) DataSource {
// 	return DataSource{
// 		ID:          id,
// 		Name:        name,
// 		Description: description,
// 		Owner:       owner.String(),
// 		Fee:         fee.String(),
// 		Executable:  executable,
// 		LastUpdated: blockTime.UnixNano() / int64(time.Millisecond),
// 	}
// }

// func (b *BandDB) AddDataSource(
// 	id int64,
// 	name, description string,
// 	owner sdk.AccAddress,
// 	fee sdk.Coins,
// 	executable []byte,
// 	blockTime time.Time,
// 	blockHeight int64,
// 	txHash []byte,
// ) error {
// 	dataSource := createDataSource(
// 		id,
// 		name,
// 		description,
// 		owner,
// 		fee,
// 		executable,
// 		blockTime,
// 	)
// 	err := b.tx.Create(&dataSource).Error
// 	if err != nil {
// 		return err
// 	}

// 	return b.tx.Create(&DataSourceRevision{
// 		DataSourceID: id,
// 		Name:         name,
// 		Timestamp:    blockTime.UnixNano() / int64(time.Millisecond),
// 		BlockHeight:  blockHeight,
// 		TxHash:       txHash,
// 	}).Error
// }

// func (b *BandDB) handleMsgCreateDataSource(
// 	txHash []byte,
// 	msg oracle.MsgCreateDataSource,
// 	events map[string]string,
// ) error {
// 	rawID, ok := events[oracle.EventTypeCreateDataSource+"."+oracle.AttributeKeyID]
// 	if !ok {
// 		return errors.New("handleMsgCreateDataSource: cannot find data source id")
// 	}
// 	id, err := strconv.ParseInt(rawID, 10, 64)
// 	if err != nil {
// 		return err
// 	}
// 	return b.AddDataSource(
// 		id, msg.Name, msg.Description, msg.Owner, msg.Fee, msg.Executable,
// 		b.ctx.BlockTime(), b.ctx.BlockHeight(), txHash,
// 	)
// }

// func (b *BandDB) handleMsgEditDataSource(
// 	txHash []byte,
// 	msg oracle.MsgEditDataSource,
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
// 		Timestamp:    b.ctx.BlockTime().UnixNano() / int64(time.Millisecond),
// 		BlockHeight:  b.ctx.BlockHeight(),
// 		TxHash:       txHash,
// 	}).Error
// }
