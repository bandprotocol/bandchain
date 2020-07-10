package yoda

import (
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"

	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
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

		if messageType == (otypes.MsgRequestData{}).Type() {
			go handleRequestLog(c, l, log)
		} else {
			l.Debug(":ghost: Skipping non-{request/packet} type: %s", messageType)
		} /*else if messageType == (ibc.MsgPacket{}).Type() {
			// Try to get request id from packet. If not then return error.
			_, err := GetEventValue(log, otypes.EventTypeRequest, otypes.AttributeKeyID)
			if err != nil {
				l.Debug(":ghost: Skipping non-request packet")
				return
			}
			go handleRequestLog(c, l, log)
		} */
	}
}

func handleRequestLog(c *Context, l *Logger, log sdk.ABCIMessageLog) {
	idStr, err := GetEventValue(log, otypes.EventTypeRequest, otypes.AttributeKeyID)
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
	validators := GetEventValues(log, otypes.EventTypeRequest, otypes.AttributeKeyValidator)
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

	reportsChan := make(chan otypes.RawReport, len(reqs))
	for _, req := range reqs {
		go func(l *Logger, req rawRequest) {
			exec, err := GetExecutable(c, l, req.dataSourceHash)
			if err != nil {
				l.Error(":skull: Failed to load data source with error: %s", err.Error())
				reportsChan <- otypes.NewRawReport(
					req.externalID, 255, []byte("FAIL_TO_LOAD_DATA_SOURCE"),
				)
				return
			}
			// TODO: Make timeout can configurable.
			result, err := c.executor.Exec(10*time.Second, exec, req.calldata)
			if err != nil {
				l.Error(":skull: Failed to execute data source script: %s", err.Error())
				reportsChan <- otypes.NewRawReport(req.externalID, 255, nil)
			} else {
				l.Debug(
					":sparkles: Query data done with calldata: %q, result: %q, exitCode: %d",
					req.calldata, result.Output, result.Code,
				)
				reportsChan <- otypes.NewRawReport(req.externalID, result.Code, result.Output)
			}
		}(l.With("did", req.dataSourceID, "eid", req.externalID), req)
	}

	reports := make([]otypes.RawReport, 0)
	for range reqs {
		reports = append(reports, <-reportsChan)
	}

	SubmitReport(c, l, otypes.RequestID(id), reports)
}
