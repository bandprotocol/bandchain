package common

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EvMap is a type alias for SDK events mapping from Attr.Key to the list of values.
type EvMap map[string][]string

// JsDict is a type alias for JSON dictionary.
type JsDict map[string]interface{}

// Message is a simple wrapper data type for each message published to Kafka.
type Message struct {
	Key   string
	Value JsDict
}

type Emitter interface {
	Write(key string, val JsDict)
}

func EmitCommit(e Emitter, height int64) {
	e.Write("COMMIT", JsDict{"height": height})
}

func EmitNewBlock(e Emitter, height int64, timestamp int64, proposer string, hash []byte, inflation string, supply string) {
	e.Write("NEW_BLOCK", JsDict{
		"height":    height,
		"timestamp": timestamp,
		"proposer":  proposer,
		"hash":      hash,
		"inflation": inflation,
		"supply":    supply,
	})
}

func EmitNewUnbondingDelegation(e Emitter, delegator sdk.AccAddress, validator sdk.ValAddress, height, completeTime int64, amount sdk.Int) {
	e.Write("NEW_UNBONDING_DELEGATION", JsDict{
		"delegator_address": delegator,
		"operator_address":  validator,
		"creation_height":   height,
		"completion_time":   completeTime,
		"amount":            amount,
	})
}

func EmitNewRedelegation(e Emitter, delegator sdk.AccAddress, valSrcAddr sdk.ValAddress, valDstAddr sdk.ValAddress, completeTime int64, amount sdk.Int) {
	e.Write("NEW_REDELEGATION", JsDict{
		"delegator_address":    delegator,
		"operator_src_address": valSrcAddr,
		"operator_dst_address": valDstAddr,
		"completion_time":      completeTime,
		"amount":               amount,
	})
}

func EmitNewValidatorVote(e Emitter, consensusAddress string, height int64, voted bool) {
	e.Write("NEW_VALIDATOR_VOTE", JsDict{
		"consensus_address": consensusAddress,
		"block_height":      height,
		"voted":             voted,
	})
}

func EmitSetRelatedTransaction(e Emitter, hash []byte, addrs []sdk.AccAddress) {
	e.Write("SET_RELATED_TRANSACTION", JsDict{
		"hash":             hash,
		"related_accounts": addrs,
	})
}

func EmitNewDataSourceRequest(e Emitter, id types.DataSourceID) {
	e.Write("NEW_DATA_SOURCE_REQUEST", JsDict{
		"data_source_id": id,
		"count":          0,
	})
}

func EmitNewOracleScriptRequest(e Emitter, id types.OracleScriptID) {
	e.Write("NEW_ORACLE_SCRIPT_REQUEST", JsDict{
		"oracle_script_id": id,
		"count":            0,
	})
}

func EmitUpdateDataSourceRequest(e Emitter, id types.DataSourceID) {
	e.Write("UPDATE_DATA_SOURCE_REQUEST", JsDict{
		"data_source_id": id,
	})
}

func EmitUpdateOracleScriptRequest(e Emitter, id types.OracleScriptID) {
	e.Write("UPDATE_ORACLE_SCRIPT_REQUEST", JsDict{
		"oracle_script_id": id,
	})
}

func EmitUpdateRelatedDsOs(e Emitter, dsID types.DataSourceID, osID types.OracleScriptID) {
	e.Write("UPDATE_RELATED_DS_OS", JsDict{
		"oracle_script_id": osID,
		"data_source_id":   dsID,
	})
}

func EmitSetHistoricalBondedTokenOnValidator(e Emitter, addr sdk.ValAddress, tokens uint64, timestamp int64) {
	e.Write("SET_HISTORICAL_BONDED_TOKEN_ON_VALIDATOR", JsDict{
		"operator_address": addr,
		"bonded_tokens":    tokens,
		"timestamp":        timestamp,
	})
}

func EmitSetRequestCountPerDay(e Emitter, date int64) {
	e.Write("SET_REQUEST_COUNT_PER_DAY", JsDict{
		"date": date,
	})
}

func EmitHistoricalValidatorStatus(e Emitter, addr sdk.ValAddress, status bool, timestamp int64) {
	e.Write("SET_HISTORICAL_VALIDATOR_STATUS", JsDict{
		"operator_address": addr,
		"status":           status,
		"timestamp":        timestamp,
	})
}
