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
	"github.com/cosmos/cosmos-sdk/x/staking"
)

type BandDB struct {
	db  *gorm.DB
	tx  *gorm.DB
	ctx sdk.Context

	StakingKeeper staking.Keeper
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
		&RelatedDataSources{},
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

	db.Model(&Delegation{}).AddForeignKey(
		"operator_address",
		"validators(operator_address)",
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

	db.Model(&RelatedDataSources{}).AddForeignKey(
		"data_source_id",
		"data_sources(id)",
		"RESTRICT",
		"RESTRICT",
	)

	db.Model(&RelatedDataSources{}).AddForeignKey(
		"oracle_script_id",
		"oracle_scripts(id)",
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

func (b *BandDB) HandleTransaction(tx auth.StdTx, txHash []byte, logs sdk.ABCIMessageLogs) {
	msgs := tx.GetMsgs()

	if len(msgs) != len(logs) {
		panic("Inconsistent size of msgs and logs.")
	}

	messages := make([]map[string]interface{}, 0)

	for idx, msg := range msgs {
		events := logs[idx].Events
		kvMap := make(map[string]string)
		for _, event := range events {
			for _, kv := range event.Attributes {
				kvMap[event.Type+"."+kv.Key] = kv.Value
			}
		}

		newMsg, err := b.HandleMessage(txHash, msg, kvMap)
		if err != nil {
			panic(err)
		}

		messages = append(messages, newMsg)

	}

	b.UpdateTransaction(txHash, messages)
}

func (b *BandDB) HandleMessage(txHash []byte, msg sdk.Msg, events map[string]string) (map[string]interface{}, error) {
	jsonMap := make(map[string]interface{})
	rawBytes, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rawBytes, &jsonMap)
	if err != nil {
		return nil, err
	}

	switch msg := msg.(type) {
	case zoracle.MsgCreateDataSource:
		err = b.handleMsgCreateDataSource(txHash, msg, events)
		if err != nil {
			return nil, err
		}

		jsonMap["dataSourceID"] = events[zoracle.EventTypeCreateDataSource+"."+zoracle.AttributeKeyID]
	case zoracle.MsgEditDataSource:
		err = b.handleMsgEditDataSource(txHash, msg, events)
		if err != nil {
			return nil, err
		}
	case zoracle.MsgCreateOracleScript:
		err = b.handleMsgCreateOracleScript(txHash, msg, events)
		if err != nil {
			return nil, err
		}

		jsonMap["oracleScriptID"] = events[zoracle.EventTypeCreateOracleScript+"."+zoracle.AttributeKeyID]
	case zoracle.MsgEditOracleScript:
		err = b.handleMsgEditOracleScript(txHash, msg, events)
		if err != nil {
			return nil, err
		}
	case zoracle.MsgRequestData:
		err = b.handleMsgRequestData(txHash, msg, events)
		if err != nil {
			return nil, err
		}

		var oracleScript OracleScript
		err := b.tx.First(&oracleScript, int64(msg.OracleScriptID)).Error
		if err != nil {
			return nil, err
		}
		jsonMap["oracleScriptName"] = oracleScript.Name
		jsonMap["requestID"] = events[zoracle.EventTypeRequest+"."+zoracle.AttributeKeyID]
	case zoracle.MsgReportData:
		err = b.handleMsgReportData(txHash, msg, events)
		if err != nil {
			return nil, err
		}
	case zoracle.MsgAddOracleAddress:
		val, _ := b.StakingKeeper.GetValidator(b.ctx, msg.Validator)
		jsonMap["validatorMoniker"] = val.Description.Moniker
	case zoracle.MsgRemoveOracleAddress:
		val, _ := b.StakingKeeper.GetValidator(b.ctx, msg.Validator)
		jsonMap["validatorMoniker"] = val.Description.Moniker
	case bank.MsgSend:
		return b.handleMsgRequestData(txHash, msg, events)
	case staking.MsgCreateValidator:
		return b.handleMsgCreateValidator(msg)
	case staking.MsgEditValidator:
		return b.handleMsgEditValidator(msg)
	default:
		// TODO: Better logging
		fmt.Println("HandleMessage: There isn't event handler for this type")
		return nil, nil
	}
	jsonMap["type"] = events["message.action"]

	return jsonMap, nil
}
