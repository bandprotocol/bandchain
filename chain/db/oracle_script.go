package db

import (
	"errors"
	"strconv"
	"time"

	"github.com/bandprotocol/bandchain/chain/x/zoracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createOracleScript(
	id int64,
	name, description string,
	owner sdk.AccAddress,
	code []byte,
	blockTime time.Time,
) OracleScript {

	return OracleScript{
		ID:          id,
		Name:        name,
		Description: description,
		Owner:       owner.String(),
		Code:        code,
		LastUpdated: blockTime.UnixNano() / int64(time.Millisecond),
	}
}

func (b *BandDB) AddOracleScript(
	id int64,
	name, description string,
	owner sdk.AccAddress,
	code []byte,
	blockTime time.Time,
	blockHeight int64,
	txHash []byte,
) error {
	// codeHash := sha256.New().Sum(code)

	oracleScript := createOracleScript(
		id,
		name,
		description,
		owner,
		code,
		blockTime,
	)
	err := b.tx.Create(&oracleScript).Error
	if err != nil {
		return err
	}

	return b.tx.Create(&OracleScriptRevision{
		OracleScriptID: id,
		Name:           name,
		Timestamp:      blockTime.UnixNano() / int64(time.Millisecond),
		BlockHeight:    blockHeight,
		TxHash:         txHash,
	}).Error
}

func (b *BandDB) handleMsgCreateOracleScript(
	txHash []byte,
	msg zoracle.MsgCreateOracleScript,
	events map[string]string,
) error {
	rawID, ok := events[zoracle.EventTypeCreateOracleScript+"."+zoracle.AttributeKeyID]
	if !ok {
		return errors.New("handleMsgCreateOracleScript: cannot find oracle script id")
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return err
	}
	return b.AddOracleScript(
		id, msg.Name, msg.Description, msg.Owner, msg.Code,
		b.ctx.BlockTime(), b.ctx.BlockHeight(), txHash,
	)
}

func (b *BandDB) handleMsgEditOracleScript(
	txHash []byte,
	msg zoracle.MsgEditOracleScript,
	events map[string]string,
) error {
	oracleScript := createOracleScript(
		int64(msg.OracleScriptID), msg.Name, msg.Description,
		msg.Owner, msg.Code, b.ctx.BlockTime(),
	)

	err := b.tx.Save(&oracleScript).Error
	if err != nil {
		return err
	}

	return b.tx.Create(&OracleScriptRevision{
		OracleScriptID: int64(msg.OracleScriptID),
		Name:           msg.Name,
		Timestamp:      b.ctx.BlockTime().UnixNano() / int64(time.Millisecond),
		BlockHeight:    b.ctx.BlockHeight(),
		TxHash:         txHash,
	}).Error
}
