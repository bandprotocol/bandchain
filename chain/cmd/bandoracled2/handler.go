package main

import (
	"sort"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bandprotocol/bandchain/chain/byteexec"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

func handleTransaction(c *Context, l *Logger, tx tmtypes.TxResult) {
	l.Debug(":eyes: Inspecting incoming transaction %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		l.Debug(":alien: Skipping transaction with non-zero code %d", tx.Result.Code)
		return
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		l.Error(":cold_sweat: Failed to parse transaction logs with error: %s", err.Error())
		return
	}

	for _, log := range logs {
		// TODO: Also handle IBC request packet here
		messageType, err := GetEventValue(log, "message", "action")
		if err != nil {
			l.Error(":cold_sweat: Failed to get message action type with error: %s", err.Error())
			continue
		}
		if messageType != "request" {
			l.Debug(":ghost: Skipping non-request message type: %s", messageType)
			continue
		}
		go handleRequestLog(c, l, log)
	}
}

func handleRequestLog(c *Context, l *Logger, log sdk.ABCIMessageLog) {
	idStr, err := GetEventValue(log, "request", "id")
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

	dataSourceIDs := GetEventValues(log, "raw_request", "data_source_id")
	externalIDs := GetEventValues(log, "raw_request", "external_id")
	calldataList := GetEventValues(log, "raw_request", "calldata")
	if len(dataSourceIDs) != len(externalIDs) {
		l.Error(":skull: Inconsistent data source count and external ID count")
		return
	}
	if len(dataSourceIDs) != len(calldataList) {
		l.Error(":skull: Inconsistent data source count and calldata count")
		return
	}

	reports := make([]oracle.RawReport, 0)
	// TODO: Parallelize the work to get data from each source.
	for idx := range dataSourceIDs {
		dataSourceID, err := strconv.Atoi(dataSourceIDs[idx])
		if err != nil {
			l.Error(":cold_sweat: Failed to parse data source id with error: %s", err.Error())
			return
		}

		externalID, err := strconv.Atoi(externalIDs[idx])
		if err != nil {
			l.Error(":cold_sweat: Failed to parse external id with error: %s", err.Error())
			return
		}

		executable, err := GetExecutable(c, dataSourceID)
		if err != nil {
			l.Error(
				":cold_sweat: Failed to get executable of data source id #%d with error: %s",
				dataSourceID, err.Error(),
			)
			return
		}

		// TODO: Parallelize the work. Remove hardcode.
		result, err := byteexec.RunOnAWSLambda(
			executable, 3*time.Second, calldataList[idx],
			"https://dmptasv4j8.execute-api.ap-southeast-1.amazonaws.com/bash-execute",
		)
		// TODO: Extract exit code.
		if err != nil {
			l.Error(
				"❌ Request #%d: failed to run executable of data source id #%d with error: %d",
				id, dataSourceID, err,
			)
			return
		}
		reports = append(reports, oracle.NewRawReport(oracle.ExternalID(externalID), 0, result))
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ExternalID < reports[j].ExternalID
	})

	BroadCastMsgs(c, []sdk.Msg{
		oracle.NewMsgReportData(
			oracle.RequestID(id), reports, sdk.ValAddress(c.key.GetAddress()), c.key.GetAddress(),
		),
	})
}
