package db

import (
	"crypto/sha256"
	"database/sql"
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
	codeHash []byte,
	blockTime time.Time,
) OracleScript {

	return OracleScript{
		ID:          id,
		Name:        name,
		Description: description,
		Owner:       owner.String(),
		CodeHash:    codeHash,
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

	h := sha256.New()
	h.Write(code)
	codeHash := h.Sum(nil)
	oracleScriptCode := OracleScriptCode{
		CodeHash: codeHash,
		CodeText: sql.NullString{},
		Schema:   sql.NullString{},
	}

	err := b.tx.Where(OracleScriptCode{CodeHash: codeHash}).
		Assign(oracleScriptCode).
		FirstOrCreate(&OracleScriptCode{}).Error
	if err != nil {
		return err
	}

	oracleScript := createOracleScript(
		id,
		name,
		description,
		owner,
		codeHash,
		blockTime,
	)
	err = b.tx.Create(&oracleScript).Error
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
	h := sha256.New()
	h.Write(msg.Code)
	codeHash := h.Sum(nil)

	err := b.tx.Where(&OracleScriptCode{CodeHash: codeHash}).
		Assign(OracleScriptCode{
			CodeText: sql.NullString{},
			Schema:   sql.NullString{},
		}).FirstOrCreate(&OracleScriptCode{}).Error

	if err != nil {
		return err
	}

	oracleScript := createOracleScript(
		int64(msg.OracleScriptID), msg.Name, msg.Description,
		msg.Owner, codeHash, b.ctx.BlockTime(),
	)

	err = b.tx.Save(&oracleScript).Error
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
