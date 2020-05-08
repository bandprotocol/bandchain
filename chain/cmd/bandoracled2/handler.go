package main

import (
	"sort"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/bandprotocol/bandchain/chain/byteexec"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

func handleTransaction(client *rpchttp.HTTP, key keyring.Info, tx tmtypes.TxResult) {
	logger.Debug("👀 Inspecting incoming transaction %X", tmhash.Sum(tx.Tx))
	if tx.Result.Code != 0 {
		logger.Debug("🤖 Skipping transaction with non-zero code %d", tx.Result.Code)
		return
	}

	logs, err := sdk.ParseABCILogs(tx.Result.Log)
	if err != nil {
		logger.Error("❌ Failed to parse transaction logs with error: %s", err.Error())
		return
	}

	for _, log := range logs {
		// TODO: Also handle IBC request packet here
		messageType, err := GetEventValue(log, "message", "action")
		if err != nil {
			logger.Error("❌ Failed to get message action type with error: %s", err.Error())
			continue
		}
		if messageType != "request" {
			logger.Debug("👻 Skipping non-request message type: %s", messageType)
			continue
		}
		go handleRequestLog(client, key, log)
	}
}

func handleRequestLog(client *rpchttp.HTTP, key keyring.Info, log sdk.ABCIMessageLog) {
	idStr, err := GetEventValue(log, "request", "id")
	if err != nil {
		logger.Error("❌ Failed to parse request id with error: %s", err.Error())
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("❌ Convert request id: %s to integer with error: %s", idStr, err.Error())
		return
	}
	logger.Info("🚚 Processing incoming request event with id: %d", id)

	dataSourceIDs := GetEventValues(log, "raw_request", "data_source_id")
	externalIDs := GetEventValues(log, "raw_request", "external_id")
	calldataList := GetEventValues(log, "raw_request", "calldata")
	if len(dataSourceIDs) != len(externalIDs) {
		logger.Error("😱 Request #%d: inconsistent data source count and external ID count", id)
		return
	}
	if len(dataSourceIDs) != len(calldataList) {
		logger.Error("😱 Request #%d: inconsistent data source count and calldata count", id)
		return
	}

	reports := make([]oracle.RawReport, 0)
	// TODO: Parallelize the work to get data from each source.
	for idx := range dataSourceIDs {
		dataSourceID, err := strconv.Atoi(dataSourceIDs[idx])
		if err != nil {
			logger.Error(
				"❌ Request #%d: failed to parse data source id %s with error: %d",
				id, dataSourceIDs[idx], err,
			)
			return
		}

		externalID, err := strconv.Atoi(externalIDs[idx])
		if err != nil {
			logger.Error(
				"❌ Request #%d: failed to parse external id %s with error: %d",
				id, externalIDs[idx], err,
			)
			return
		}

		executable, err := GetExecutable(client, dataSourceID)
		if err != nil {
			logger.Error(
				"❌ Request #%d: failed to get executable of data source id #%d with error: %d",
				id, dataSourceID, err,
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
			logger.Error(
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

	BroadCastMsgs(client, key, []sdk.Msg{
		oracle.NewMsgReportData(
			oracle.RequestID(id), reports, sdk.ValAddress(key.GetAddress()), key.GetAddress(),
		),
	})
}
