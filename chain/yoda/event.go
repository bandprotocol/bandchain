package yoda

import (
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

type rawRequest struct {
	dataSourceID   types.DataSourceID
	dataSourceHash string
	externalID     types.ExternalID
	calldata       string
}

// GetRawRequests returns the list of all raw data requests in the given log.
func GetRawRequests(log sdk.ABCIMessageLog) ([]rawRequest, error) {
	dataSourceIDs := GetEventValues(log, types.EventTypeRawRequest, types.AttributeKeyDataSourceID)
	dataSourceHashList := GetEventValues(log, types.EventTypeRawRequest, types.AttributeKeyDataSourceHash)
	externalIDs := GetEventValues(log, types.EventTypeRawRequest, types.AttributeKeyExternalID)
	calldataList := GetEventValues(log, types.EventTypeRawRequest, types.AttributeKeyCalldata)

	if len(dataSourceIDs) != len(externalIDs) {
		return nil, fmt.Errorf("Inconsistent data source count and external ID count")
	}
	if len(dataSourceIDs) != len(calldataList) {
		return nil, fmt.Errorf("Inconsistent data source count and calldata count")
	}

	var reqs []rawRequest
	for idx := range dataSourceIDs {
		dataSourceID, err := strconv.Atoi(dataSourceIDs[idx])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse data source id: %s", err.Error())
		}

		externalID, err := strconv.Atoi(externalIDs[idx])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse external id: %s", err.Error())
		}

		reqs = append(reqs, rawRequest{
			dataSourceID:   types.DataSourceID(dataSourceID),
			dataSourceHash: dataSourceHashList[idx],
			externalID:     types.ExternalID(externalID),
			calldata:       calldataList[idx],
		})
	}
	return reqs, nil
}

// GetEventValues returns the list of all values in the given log with the given type and key.
func GetEventValues(log sdk.ABCIMessageLog, evType string, evKey string) (res []string) {
	for _, ev := range log.Events {
		if ev.Type != evType {
			continue
		}

		for _, attr := range ev.Attributes {
			if attr.Key == evKey {
				res = append(res, attr.Value)
			}
		}
	}
	return res
}

// GetEventValue checks and returns the exact value in the given log with the given type and key.
func GetEventValue(log sdk.ABCIMessageLog, evType string, evKey string) (string, error) {
	values := GetEventValues(log, evType, evKey)
	if len(values) == 0 {
		return "", fmt.Errorf("Cannot find event with type: %s, key: %s", evType, evKey)
	}
	if len(values) > 1 {
		return "", fmt.Errorf("Found more than one event with type: %s, key: %s", evType, evKey)
	}
	return values[0], nil
}
