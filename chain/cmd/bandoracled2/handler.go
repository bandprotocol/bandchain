package main

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/bandprotocol/bandchain/chain/byteexec"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

func handleRequestLog(log sdk.ABCIMessageLog) error {
	// TODO: Remove magic string
	idStr, err := GetEventValue(log, "request", "id")
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}

	// TODO: Remove magic string
	dataSourceIDs := GetEventValues(log, "raw_request", "data_source_id")
	externalIDs := GetEventValues(log, "raw_request", "external_id")
	calldataList := GetEventValues(log, "raw_request", "calldata")

	if len(dataSourceIDs) != len(externalIDs) {
		return fmt.Errorf("Inconsistent count: data source id count: %d, external id count: %d",
			len(dataSourceIDs), len(externalIDs))
	}

	if len(dataSourceIDs) != len(calldataList) {
		return fmt.Errorf("Inconsistent count: data source id count: %d, calldata count: %d",
			len(dataSourceIDs), len(calldataList))
	}

	reports := make([]oracle.RawReport, 0)
	for idx, _ := range dataSourceIDs {
		dataSourceID, err := strconv.Atoi(dataSourceIDs[idx])
		if err != nil {
			return err
		}

		externalID, err := strconv.Atoi(externalIDs[idx])
		if err != nil {
			return err
		}

		executable, err := GetExecutable(dataSourceID)
		if err != nil {
			return err
		}

		// TODO: Parallelize the work. Remove hardcode.
		result, err := byteexec.RunOnAWSLambda(
			executable, 3*time.Second, calldataList[idx],
			"https://dmptasv4j8.execute-api.ap-southeast-1.amazonaws.com/bash-execute",
		)
		// TODO: Extract exit code.
		if err != nil {
			return err
		}
		reports = append(reports, oracle.NewRawReport(oracle.ExternalID(externalID), 0, result))
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ExternalID < reports[j].ExternalID
	})

	// TODO: Remove hardcode.
	tx, err := bandClient.SendTransaction(
		oracle.NewMsgReportData(
			oracle.RequestID(id), reports, sdk.ValAddress(bandClient.Sender()),
			bandClient.Sender(),
		),
		300000, "", "",
	)

	if err != nil {
		return err
	}

	logger.Info("Report on request #%d successfully. Tx: %v", id, tx)
	return nil
}
