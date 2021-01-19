package yoda

import (
	"encoding/hex"
	"strconv"

	ckeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type processingResult struct {
	rawReport types.RawReport
	version   string
	err       error
}

func handleTransaction(c *Context, l *Logger, tx tmtypes.TxResult) {
	l.Debug(":eyes: Inspecting incoming transaction: %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		l.Debug(":alien: Skipping transaction with non-zero code: %d", tx.Result.Code)
		return
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		l.Error(":cold_sweat: Failed to parse transaction logs with error: %s", c, err.Error())
		return
	}

	for _, log := range logs {
		messageType, err := GetEventValue(log, sdk.EventTypeMessage, sdk.AttributeKeyAction)
		if err != nil {
			l.Error(":cold_sweat: Failed to get message action type with error: %s", c, err.Error())
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
		l.Error(":cold_sweat: Failed to parse request id with error: %s", c, err.Error())
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		l.Error(":cold_sweat: Failed to convert %s to integer with error: %s", c, idStr, err.Error())
		return
	}

	l = l.With("rid", id)

	// If id is in pending requests list, then skip it.
	if c.pendingRequests[types.RequestID(id)] {
		l.Debug(":eyes: Request is in pending list, then skip")
		return
	}

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
		l.Error(":skull: Failed to parse raw requests with error: %s", c, err.Error())
	}

	keyIndex := c.nextKeyIndex()
	key := c.keys[keyIndex]

	reports, execVersions := handleRawRequests(c, l, types.RequestID(id), reqs, key)

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
		l.Error(":skull: Fail to parse call data: %s", c, err.Error())
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
			rawRequests: reqs,
			clientID:    clientID,
		},
	}
}

func handlePendingRequest(c *Context, l *Logger, id types.RequestID) {

	req, err := GetRequest(c, l, id)
	if err != nil {
		l.Error(":skull: Failed to get request with error: %s", c, err.Error())
		return
	}

	l.Info(":delivery_truck: Processing pending request")

	keyIndex := c.nextKeyIndex()
	key := c.keys[keyIndex]

	var rawRequests []rawRequest

	// prepare raw requests
	for _, raw := range req.RawRequests {

		hash, err := GetDataSourceHash(c, l, raw.DataSourceID)
		if err != nil {
			l.Error(":skull: Failed to get data source hash with error: %s", c, err.Error())
			return
		}

		rawRequests = append(rawRequests, rawRequest{
			dataSourceID:   raw.DataSourceID,
			dataSourceHash: hash,
			externalID:     raw.ExternalID,
			calldata:       string(raw.Calldata),
		})
	}

	// process raw requests
	reports, execVersions := handleRawRequests(c, l, id, rawRequests, key)

	c.pendingMsgs <- ReportMsgWithKey{
		msg:         types.NewMsgReportData(types.RequestID(id), reports, c.validator, key.GetAddress()),
		execVersion: execVersions,
		keyIndex:    keyIndex,
		feeEstimationData: FeeEstimationData{
			askCount:    int64(len(req.RequestedValidators)),
			minCount:    int64(req.MinCount),
			callData:    req.Calldata,
			rawRequests: rawRequests,
			clientID:    req.ClientID,
		},
	}
}

func handleRawRequests(c *Context, l *Logger, id types.RequestID, reqs []rawRequest, key keys.Info) (reports []types.RawReport, execVersions []string) {
	resultsChan := make(chan processingResult, len(reqs))
	for _, req := range reqs {
		go handleRawRequest(c, l.With("did", req.dataSourceID, "eid", req.externalID), req, key, types.RequestID(id), resultsChan)
	}

	versions := map[string]bool{}
	for range reqs {
		result := <-resultsChan
		reports = append(reports, result.rawReport)

		if result.err == nil {
			versions[result.version] = true
		}
	}

	for version := range versions {
		execVersions = append(execVersions, version)
	}

	return
}

func handleRawRequest(c *Context, l *Logger, req rawRequest, key keys.Info, id types.RequestID, processingResultCh chan processingResult) {
	c.updateHandlingGauge(1)

	exec, err := GetExecutable(c, l, req.dataSourceHash)
	if err != nil {
		l.Error(":skull: Failed to load data source with error: %s", c, err.Error())
		processingResultCh <- processingResult{
			rawReport: types.NewRawReport(
				req.externalID, 255, []byte("FAIL_TO_LOAD_DATA_SOURCE"),
			),
			err: err,
		}
		c.updateHandlingGauge(-1)
		return
	}

	vmsg := NewVerificationMessage(cfg.ChainID, c.validator, id, req.externalID)
	sig, pubkey, err := keybase.Sign(key.GetName(), ckeys.DefaultKeyPass, vmsg.GetSignBytes())
	if err != nil {
		l.Error(":skull: Failed to sign verify message: %s", c, err.Error())
		processingResultCh <- processingResult{
			rawReport: types.NewRawReport(req.externalID, 255, nil),
			err:       err,
		}
		c.updateHandlingGauge(-1)
		return
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
		l.Error(":skull: Failed to execute data source script: %s", c, err.Error())
		processingResultCh <- processingResult{
			rawReport: types.NewRawReport(req.externalID, 255, nil),
			err:       err,
		}
		c.updateHandlingGauge(-1)
		return
	} else {
		l.Debug(
			":sparkles: Query data done with calldata: %q, result: %q, exitCode: %d",
			req.calldata, result.Output, result.Code,
		)
		processingResultCh <- processingResult{
			rawReport: types.NewRawReport(req.externalID, result.Code, result.Output),
			version:   result.Version,
		}
	}
	c.updateHandlingGauge(-1)
}
