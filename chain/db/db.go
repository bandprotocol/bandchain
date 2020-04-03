package db

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/bandprotocol/bandchain/chain/x/zoracle"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	dist "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/slashing"
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
		&Block{},
		&Transaction{},
		&Account{},
		&Validator{},
		&ValidatorVote{},
		&DataSource{},
		&DataSourceRevision{},
		&OracleScript{},
		&OracleScriptCode{},
		&OracleScriptRevision{},
		&RelatedDataSources{},
		&Request{},
		&RequestedValidator{},
		&RawDataRequests{},
		&Report{},
		&ReportDetail{},
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

	db.Model(&Transaction{}).AddForeignKey(
		"sender",
		"accounts(address)",
		"RESTRICT",
		"RESTRICT",
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

	db.Model(&OracleScript{}).AddForeignKey(
		"code_hash",
		"oracle_script_codes(code_hash)",
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
		"request_id",
		"requests(id)",
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

	db.Model(&ReportDetail{}).AddForeignKey(
		"request_id,validator",
		"reports(request_id,validator)",
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

func wrapMessage(msg []map[string]interface{}, status string) map[string]interface{} {
	objMsg := make(map[string]interface{})
	objMsg["messages"] = msg
	objMsg["status"] = status
	return objMsg
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
	wrapedMsg := wrapMessage(messages, "success")
	b.UpdateTransaction(txHash, wrapedMsg)
}

func (b *BandDB) HandleTransactionFail(tx auth.StdTx, txHash []byte) {
	txMsgs := tx.GetMsgs()
	messages := make([]map[string]interface{}, 0)
	for _, txMsg := range txMsgs {
		message := make(map[string]interface{})
		message["sender"] = txMsg.GetSigners()[0].String()
		message["type"] = txMsg.Type()
		messages = append(messages, message)
	}

	wrapedMsg := wrapMessage(messages, "failure")
	b.UpdateTransaction(txHash, wrapedMsg)
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

		dataSourceID, err := strconv.ParseInt(events[zoracle.EventTypeCreateDataSource+"."+zoracle.AttributeKeyID], 10, 64)
		if err != nil {
			return nil, err
		}
		jsonMap["dataSourceID"] = dataSourceID
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
		oracleScriptID, err := strconv.ParseInt(events[zoracle.EventTypeCreateOracleScript+"."+zoracle.AttributeKeyID], 10, 64)
		if err != nil {
			return nil, err
		}
		jsonMap["oracleScriptID"] = oracleScriptID
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

		requestID, err := strconv.ParseInt(events[zoracle.EventTypeRequest+"."+zoracle.AttributeKeyID], 10, 64)
		if err != nil {
			return nil, err
		}

		jsonMap["oracleScriptName"] = oracleScript.Name
		jsonMap["requestID"] = requestID
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
	case staking.MsgCreateValidator:
		err := b.handleMsgCreateValidator(msg)
		if err != nil {
			return nil, err
		}
	case staking.MsgEditValidator:
		err := b.handleMsgEditValidator(msg)
		if err != nil {
			return nil, err
		}
	default:
		// TODO: Better logging
		fmt.Println("HandleMessage: There isn't event handler for this type")
		return nil, nil
	}
	jsonMap["type"] = events["message.action"]

	return jsonMap, nil
}

func (b *BandDB) GetInvolvedAccountsFromTx(tx auth.StdTx) []sdk.AccAddress {
	involvedAccounts := make([]sdk.AccAddress, 0)
	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case zoracle.MsgCreateDataSource:
			continue
		case zoracle.MsgEditDataSource:
			continue
		case zoracle.MsgCreateOracleScript:
			continue
		case zoracle.MsgEditOracleScript:
			continue
		case zoracle.MsgAddOracleAddress:
			continue
		case zoracle.MsgRemoveOracleAddress:
			continue
		case zoracle.MsgRequestData:
			involvedAccounts = append(involvedAccounts, msg.Sender)
		case zoracle.MsgReportData:
			involvedAccounts = append(involvedAccounts, msg.Reporter)
		case bank.MsgSend:
			involvedAccounts = append(involvedAccounts, msg.FromAddress, msg.ToAddress)
		case bank.MsgMultiSend:
			for _, input := range msg.Inputs {
				involvedAccounts = append(involvedAccounts, input.Address)
			}
			for _, output := range msg.Outputs {
				involvedAccounts = append(involvedAccounts, output.Address)
			}
		case staking.MsgCreateValidator:
			involvedAccounts = append(involvedAccounts, msg.DelegatorAddress)
		case staking.MsgEditValidator:
			continue
		case staking.MsgDelegate:
			involvedAccounts = append(involvedAccounts, msg.DelegatorAddress)
		case staking.MsgBeginRedelegate:
			involvedAccounts = append(involvedAccounts, msg.DelegatorAddress)
		case staking.MsgUndelegate:
			involvedAccounts = append(involvedAccounts, msg.DelegatorAddress)
		case dist.MsgSetWithdrawAddress:
			continue
		case dist.MsgWithdrawDelegatorReward:
			involvedAccounts = append(involvedAccounts, msg.DelegatorAddress)
		case dist.MsgWithdrawValidatorCommission:
			involvedAccounts = append(involvedAccounts, sdk.AccAddress(msg.ValidatorAddress))
		case gov.MsgDeposit:
			involvedAccounts = append(involvedAccounts, msg.Depositor)
		case gov.MsgSubmitProposal:
			involvedAccounts = append(involvedAccounts, msg.Proposer)
		case gov.MsgVote:
			continue
		case crisis.MsgVerifyInvariant:
			continue
		case slashing.MsgUnjail:
			continue
		default:
			panic(fmt.Sprintf("Message %s does not support", msg.Type()))
		}
	}
	return involvedAccounts
}

func (b *BandDB) GetInvolvedAccountsFromTransferEvents(logs sdk.ABCIMessageLogs) []sdk.AccAddress {
	involvedAccounts := make([]sdk.AccAddress, 0)
	for _, log := range logs {
		for _, event := range log.Events {
			if event.Type == bank.EventTypeTransfer {
				for _, kv := range event.Attributes {
					if kv.Key == bank.AttributeKeySender || kv.Key == bank.AttributeKeyRecipient {
						account, err := sdk.AccAddressFromBech32(kv.Value)
						if err != nil {
							panic(err)
						}
						involvedAccounts = append(involvedAccounts, account)
					}
				}
			}
		}
	}
	return involvedAccounts
}

func (b *BandDB) ResolveRequest(id int64, resolveStatus zoracle.ResolveStatus, result []byte) error {
	if resolveStatus == 1 {
		return b.tx.Model(&Request{}).Where(Request{ID: id}).
			Update(Request{ResolveStatus: parseResolveStatus(resolveStatus), Result: result}).Error
	}
	return b.tx.Model(&Request{}).Where(Request{ID: id}).
		Update(Request{ResolveStatus: parseResolveStatus(resolveStatus)}).Error
}
