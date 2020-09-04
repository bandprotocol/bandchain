package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

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

func EmitNewUnbondingDelegation(e Emitter, delegator sdk.AccAddress, validator sdk.ValAddress, completeTime int64, amount sdk.Int) {
	e.Write("NEW_UNBONDING_DELEGATION", JsDict{
		"delegator_address": delegator,
		"operator_address":  validator,
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
