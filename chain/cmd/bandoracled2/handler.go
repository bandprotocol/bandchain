package main

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
		// TODO: Also handle IBC request packet here
		messageType, err := GetEventValue(log, otypes.EventTypeMessage, otypes.AttributeKeyAction)
		if err != nil {
			l.Error(":cold_sweat: Failed to get message action type with error: %s", err.Error())
			continue
		}

		if messageType == otypes.MessageTypeRequest {
			go handleRequestLog(c, l, log)
		} else if messageType == otypes.MessageTypeOraclePacket {
			// Try to get request id from packet. If not then return error.
			_, err := GetEventValue(log, otypes.EventTypeRequest, otypes.AttributeKeyID)
			if err != nil {
				l.Error(":cold_sweat: Failed to parse request id from oracle packet with error: %s", err.Error())
				return
			}
			go handleRequestLog(c, l, log)
		} else {
			l.Debug(":ghost: Skipping non-request message type: %s", messageType)
		}
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
	l.Info(":delivery_truck: Processing incoming request event")

	// TODO: Skip if not related to this validator
	reqs, err := GetRawRequests(log)
	if err != nil {
		l.Error(":skull: Failed to parse raw requests with error: %s", err.Error())
	}

	reportsChan := make(chan otypes.RawReport, len(reqs))
	for _, req := range reqs {
		go func(l *Logger, req rawRequest) {
			exec, err := GetExecutable(c, l, int(req.dataSourceID))
			if err != nil {
				l.Error(":skull: Failed to load data source with error: %s", err.Error())
				reportsChan <- otypes.NewRawReport(
					req.externalID, 255, []byte("FAIL_TO_LOAD_DATA_SOURCE"),
				)
				return
			}
			// TODO: Allow user to configure different executors
			executor := &lambdaExecutor{}
			result, exitCode := executor.Execute(l, exec, 3*time.Second, req.calldata)
			l.Debug(
				":sparkles: Query data done with calldata: %q, result: %q, exitCode: %d",
				req.calldata, result, exitCode,
			)
			reportsChan <- otypes.NewRawReport(req.externalID, exitCode, result)
		}(l.With("did", req.dataSourceID, "eid", req.externalID), req)
	}

	reports := make([]otypes.RawReport, 0)
	for range reqs {
		reports = append(reports, <-reportsChan)
	}

	SubmitReport(c, l, otypes.RequestID(id), reports)
}
