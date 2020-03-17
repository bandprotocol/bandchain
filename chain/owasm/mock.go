package owasm

import (
	"fmt"
)

type mockExecutionEnvironment struct {
	requestID                         int64
	requestedValidatorCount           int64
	sufficientValidatorCount          int64
	receivedValidatorCount            int64
	prepareBlockTime                  int64
	aggregateBlockTime                int64
	validatorAddresses                [][]byte
	externalDataResults               [][][]byte
	maximumResultSize                 int64
	maximumCalldataOfDataSourceSize   int64
	requestExternalDataResultsCounter [][]int64
}

func (m *mockExecutionEnvironment) GetCurrentRequestID() int64 {
	return m.requestID
}

func (m *mockExecutionEnvironment) GetRequestedValidatorCount() int64 {
	return m.requestedValidatorCount
}

func (m *mockExecutionEnvironment) GetSufficientValidatorCount() int64 {
	return m.sufficientValidatorCount
}

func (m *mockExecutionEnvironment) GetReceivedValidatorCount() int64 {
	return m.receivedValidatorCount
}

func (m *mockExecutionEnvironment) GetPrepareBlockTime() int64 {
	return m.prepareBlockTime
}

func (m *mockExecutionEnvironment) GetAggregateBlockTime() int64 {
	return m.aggregateBlockTime
}

func (m *mockExecutionEnvironment) GetValidatorAddress(validatorIndex int64) ([]byte, error) {
	return m.validatorAddresses[validatorIndex], nil
}

func (m *mockExecutionEnvironment) GetMaximumResultSize() int64 {
	return m.maximumResultSize
}

func (m *mockExecutionEnvironment) GetMaximumCalldataOfDataSourceSize() int64 {
	return m.maximumCalldataOfDataSourceSize
}

func (m *mockExecutionEnvironment) RequestExternalData(
	dataSourceID int64,
	externalDataID int64,
	calldata []byte,
) error {
	// TODO: Figure out how to mock this elegantly.
	fmt.Printf("RequestExternalData: DataSourceID = %d, ExternalDataID = %d\n", dataSourceID, externalDataID)
	return nil
}

func (m *mockExecutionEnvironment) GetExternalData(
	externalDataID int64,
	validatorIndex int64,
) ([]byte, uint8, error) {
	if len(m.requestExternalDataResultsCounter) <= int(externalDataID) {
		return nil, 0, fmt.Errorf("externalDataID is out of range")
	}

	if len(m.requestExternalDataResultsCounter[externalDataID]) <= int(validatorIndex) {
		return nil, 0, fmt.Errorf("validatorIndex is out of range")
	}

	m.requestExternalDataResultsCounter[externalDataID][validatorIndex]++
	return m.externalDataResults[externalDataID][validatorIndex], 0, nil
}
