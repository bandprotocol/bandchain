package yoda

import (
	"encoding/hex"
	"strconv"
	"sync"

	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func handleTransaction(c *Context, l *Logger, tx tmtypes.TxResult) {
	l.Debug(":eyes: Inspecting incoming transaction: %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		l.Debug(":alien: Skipping transaction with non-zero code: %d", tx.Result.Code)
		return
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		l.Error(":cold_sweat: Failed to parse transaction logs with error: %s", err.Error())
		return
	}

	for _, log := range logs {
		messageType, err := GetEventValue(log, sdk.EventTypeMessage, sdk.AttributeKeyAction)
		if err != nil {
			l.Error(":cold_sweat: Failed to get message action type with error: %s", err.Error())
			continue
		}

		if messageType == (types.MsgRequestData{}).Type() {
			go handleRequestLog(c, l, log)
		} else {
			l.Debug(":ghost: Skipping non-{request/packet} type: %s", messageType)
		} /*else if messageType == (ibc.MsgPacket{}).Type() {
			// Try to get request id from packet. If not then return error.
			_, err := GetEventValue(log, types.EventTypeRequest, types.AttributeKeyID)
			if err != nil {
				l.Debug(":ghost: Skipping non-request packet")
				return
			}
			go handleRequestLog(c, l, log)
		} */
	}
}

func handleRequestLog(c *Context, l *Logger, log sdk.ABCIMessageLog) {
	idStr, err := GetEventValue(log, types.EventTypeRequest, types.AttributeKeyID)
	if err != nil {
		l.Error(":cold_sweat: Failed to parse request id with error: %s", err.Error())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Error(":cold_sweat: Failed to convert %s to integer with error: %s", idStr, err.Error())
		return
	}

	l = l.With("rid", id)

	// Skip if not related to this validator
	validators := GetEventValues(log, types.EventTypeRequest, types.AttributeKeyValidator)
	hasMe := false
	for _, validator := range validators {
		if validator == c.validator.String() {
			hasMe = true
			break
		}
	}

	if !hasMe {
		l.Debug(":next_track_button: Skip request not related to this validator")
		return
	}

	l.Info(":delivery_truck: Processing incoming request event")

	reqs, err := GetRawRequests(log)
	if err != nil {
		l.Error(":skull: Failed to parse raw requests with error: %s", err.Error())
	}

	keyIndex := c.nextKeyIndex()
	key := c.keys[keyIndex]

	reportsChan := make(chan types.RawReport, len(reqs))
	var version sync.Map
	for _, req := range reqs {
		go func(l *Logger, req rawRequest) {
			exec, err := GetExecutable(c, l, req.dataSourceHash)
			if err != nil {
				l.Error(":skull: Failed to load data source with error: %s", err.Error())
				reportsChan <- types.NewRawReport(
					req.externalID, 255, []byte("FAIL_TO_LOAD_DATA_SOURCE"),
				)
				return
			}

			vmsg := NewVerificationMessage(cfg.ChainID, c.validator, types.RequestID(id), req.externalID)
			sig, pubkey, err := keybase.Sign(key.GetName(), ckeys.DefaultKeyPass, vmsg.GetSignBytes())
			if err != nil {
				l.Error(":skull: Failed to sign verify message: %s", err.Error())
				reportsChan <- types.NewRawReport(req.externalID, 255, nil)
			}

			result, err := c.executor.Exec(exec, req.calldata, map[string]interface{}{
				"BAND_CHAIN_ID":    vmsg.ChainID,
				"BAND_VALIDATOR":   vmsg.Validator.String(),
				"BAND_REQUEST_ID":  strconv.Itoa(int(vmsg.RequestID)),
				"BAND_EXTERNAL_ID": strconv.Itoa(int(vmsg.ExternalID)),
				"BAND_REPORTER":    sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeAccPub, pubkey),
				"BAND_SIGNATURE":   sig,
			})

			if err != nil {
				l.Error(":skull: Failed to execute data source script: %s", err.Error())
				reportsChan <- types.NewRawReport(req.externalID, 255, nil)
			} else {
				l.Debug(
					":sparkles: Query data done with calldata: %q, result: %q, exitCode: %d",
					req.calldata, result.Output, result.Code,
				)
				version.Store(result.Version, true)
				reportsChan <- types.NewRawReport(req.externalID, result.Code, result.Output)
			}
		}(l.With("did", req.dataSourceID, "eid", req.externalID), req)
	}

	reports := make([]types.RawReport, 0)
	execVersions := make([]string, 0)
	for range reqs {
		reports = append(reports, <-reportsChan)
	}
	version.Range(func(key, value interface{}) bool {
		execVersions = append(execVersions, key.(string))
		return true
	})

	rawAskCount := GetEventValues(log, types.EventTypeRequest, types.AttributeKeyAskCount)
	if len(rawAskCount) != 1 {
		panic("Fail to get ask count")
	}
	askCount := common.Atoi(rawAskCount[0])

	rawMinCount := GetEventValues(log, types.EventTypeRequest, types.AttributeKeyMinCount)
	if len(rawMinCount) != 1 {
		panic("Fail to get min count")
	}
	minCount := common.Atoi(rawMinCount[0])

	rawCallData := GetEventValues(log, types.EventTypeRequest, types.AttributeKeyCalldata)
	if len(rawCallData) != 1 {
		panic("Fail to get call data")
	}
	callData, err := hex.DecodeString(rawCallData[0])
	if err != nil {
		l.Error(":skull: Fail to parse call data: %s", err.Error())
	}

	var clientID string
	rawClientID := GetEventValues(log, types.EventTypeRequest, types.AttributeKeyClientID)
	if len(rawClientID) > 0 {
		clientID = rawClientID[0]
	}

	c.pendingMsgs <- ReportMsgWithKey{
		msg:         types.NewMsgReportData(types.RequestID(id), reports, c.validator, key.GetAddress()),
		execVersion: execVersions,
		keyIndex:    keyIndex,
		feeEstimationData: FeeEstimationData{
			askCount:    askCount,
			minCount:    minCount,
			callData:    callData,
			validators:  len(validators),
			rawRequests: reqs,
			clientID:    clientID,
		},
	}
}
