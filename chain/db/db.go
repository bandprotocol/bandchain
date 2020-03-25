package db

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/bandprotocol/bandchain/chain/x/zoracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
)

type BandDB struct {
	db  *gorm.DB
	tx  *gorm.DB
	ctx sdk.Context

	ZoracleKeeper zoracle.Keeper
}

func NewDB(dialect, path string, metadata map[string]string) (*BandDB, error) {
	db, err := gorm.Open(dialect, path)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(
		&Metadata{},
		&Event{},
		&Validator{},
		&ValidatorVote{},
		&DataSource{},
		&DataSourceRevision{},
		&OracleScript{},
		&OracleScriptRevision{},
		&Block{},
		&Transaction{},
		&Report{},
		&ReportDetail{},
		&Request{},
		&RequestedValidator{},
		&RawDataRequests{},
	)

	db.Model(&ValidatorVote{}).AddForeignKey(
		"consensus_address",
		"validators(consensus_address)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&DataSourceRevision{}).AddForeignKey(
		"data_source_id",
		"data_sources(id)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&DataSourceRevision{}).AddForeignKey(
		"tx_hash",
		"transactions(tx_hash)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&OracleScriptRevision{}).AddForeignKey(
		"oracle_script_id",
		"oracle_scripts(id)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&OracleScriptRevision{}).AddForeignKey(
		"tx_hash",
		"transactions(tx_hash)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&Report{}).AddForeignKey(
		"validator",
		"validators(operator_address)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&Report{}).AddForeignKey(
		"tx_hash",
		"transactions(tx_hash)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&Report{}).AddForeignKey(
		"request_id",
		"requests(id)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&ReportDetail{}).AddForeignKey(
		"request_id,validator",
		"reports(request_id,validator)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&Block{}).AddForeignKey(
		"proposer",
		"validators(consensus_address)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&Transaction{}).AddForeignKey(
		"block_height",
		"blocks(height)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&Request{}).AddForeignKey(
		"oracle_script_id",
		"oracle_scripts(id)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&Request{}).AddForeignKey(
		"tx_hash",
		"transactions(tx_hash)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&RequestedValidator{}).AddForeignKey(
		"request_id",
		"requests(id)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&RequestedValidator{}).AddForeignKey(
		"validator_address",
		"validators(operator_address)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&RawDataRequests{}).AddForeignKey(
		"request_id",
		"requests(id)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&RawDataRequests{}).AddForeignKey(
		"data_source_id",
		"data_sources(id)",
		"RESTRICT",
		"RESTRICT",
	)

	for key, value := range metadata {
		err := db.Where(Metadata{Key: key}).
			Assign(Metadata{Value: value}).
			FirstOrCreate(&Metadata{}).Error
		if err != nil {
			panic(err)
		}
	}

	return &BandDB{db: db}, nil
}

func (b *BandDB) BeginTransaction() {
	if b.tx != nil {
		panic("BeginTransaction: Cannot begin a new transaction without closing the pending one.")
	}
	b.tx = b.db.Begin()
	if b.tx.Error != nil {
		panic(b.tx.Error)
	}
}

func (b *BandDB) Commit() {
	err := b.tx.Commit().Error
	if err != nil {
		panic(err)
	}
	b.tx = nil
}

func (b *BandDB) RollBack() {
	err := b.tx.Rollback()
	if err != nil {
		panic(err)
	}
	b.tx = nil
}

func (b *BandDB) SetContext(ctx sdk.Context) {
	b.ctx = ctx
}

func (b *BandDB) HandleTransaction(tx auth.StdTx, txHash []byte, logs sdk.ABCIMessageLogs) string {
	msgs := tx.GetMsgs()

	if len(msgs) != len(logs) {
		panic("Inconsistent size of msgs and logs.")
	}

	messages := []string{}

	for idx, msg := range msgs {
		events := logs[idx].Events
		kvMap := make(map[string]string)
		for _, event := range events {
			for _, kv := range event.Attributes {
				kvMap[event.Type+"."+kv.Key] = kv.Value
			}
		}

		str, err := b.HandleMessage(txHash, msg, kvMap)
		if err != nil {
			panic(err)
		}

		messages = append(messages, str)

	}
	return fmt.Sprint(messages)
}

func (b *BandDB) HandleMessage(txHash []byte, msg sdk.Msg, events map[string]string) (string, error) {
	jsonMap := make(map[string]interface{})
	jsonStr, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(jsonStr), &jsonMap)
	if err != nil {
		return "", err
	}
	switch msg := msg.(type) {
	// Just proof of concept
	case zoracle.MsgCreateDataSource:

		err = b.handleMsgCreateDataSource(txHash, msg, events)
		if err != nil {
			return "", err
		}

		jsonMap["dataSourceID"] = events[zoracle.EventTypeCreateDataSource+"."+zoracle.AttributeKeyID]

	case zoracle.MsgEditDataSource:

		err = b.handleMsgEditDataSource(txHash, msg, events)
		if err != nil {
			return "", err
		}

	case zoracle.MsgCreateOracleScript:

		jsonMap["oracleScriptID"] = events[zoracle.EventTypeCreateOracleScript+"."+zoracle.AttributeKeyID]

	case zoracle.MsgEditOracleScript:

	case zoracle.MsgRequestData:
		err := b.handleMsgRequestData(txHash, msg, events)
		if err != nil {
			return "", err
		}
		jsonMap[zoracle.AttributeKeyRequestID] = events["request.id"]
		jsonMap["type"] = events["message.action"]

		jsonMap["requestID"] = events[zoracle.EventTypeRequest+"."+zoracle.AttributeKeyID]

	case zoracle.MsgReportData:

	case zoracle.MsgAddOracleAddress:

	case zoracle.MsgRemoveOracleAdderess:

	case bank.MsgSend:

	default:
		// TODO: Better logging
		fmt.Println("HandleMessage: There isn't event handler for this type")
		return "", nil
	}
	jsonMap["type"] = events["message.action"]

	jsonString, err := json.Marshal(jsonMap)
	if err != nil {
		return "", err
	}
	return fmt.Sprint(string(jsonString)), nil

}
