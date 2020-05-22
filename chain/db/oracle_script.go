package db

import (
	"errors"
	"strconv"
	"time"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func createOracleScript(
	id int64,
	name, description string,
	owner sdk.AccAddress,
	blockTime time.Time,
	schema string,
	sourceCodeURL string,
) OracleScript {

	return OracleScript{
		ID:            id,
		Name:          name,
		Description:   description,
		Owner:         owner.String(),
		LastUpdated:   blockTime.UnixNano() / int64(time.Millisecond),
		Schema:        schema,
		SourceCodeURL: sourceCodeURL,
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
	schema string,
	sourceCodeURL string,
) error {

	oracleScript := createOracleScript(
		id,
		name,
		description,
		owner,
		blockTime,
		schema,
		sourceCodeURL,
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
	msg oracle.MsgCreateOracleScript,
	events map[string]interface{},
) error {
	rawID, ok := events[otypes.EventTypeCreateOracleScript+"."+otypes.AttributeKeyID].(string)
	if !ok {
		return errors.New("handleMsgCreateOracleScript: cannot find oracle script id")
	}
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		return err
	}
	return b.AddOracleScript(
		id, msg.Name, msg.Description, msg.Owner, msg.Code,
		b.ctx.BlockTime(), b.ctx.BlockHeight(), txHash, msg.Schema, msg.SourceCodeURL,
	)
}

func (b *BandDB) handleMsgEditOracleScript(
	txHash []byte,
	msg oracle.MsgEditOracleScript,
) error {
	oracleScript := createOracleScript(
		int64(msg.OracleScriptID), msg.Name, msg.Description,
		msg.Owner, b.ctx.BlockTime(), msg.Schema, msg.SourceCodeURL,
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
