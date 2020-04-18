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
	externalID int64,
	calldata []byte,
) error {
	// TODO: Figure out how to mock this elegantly.
	fmt.Printf("RequestExternalData: DataSourceID = %d, ExternalID = %d\n", dataSourceID, externalID)
	return nil
}

func (m *mockExecutionEnvironment) GetExternalData(
	externalID int64,
	validatorIndex int64,
) ([]byte, uint8, error) {
	if int64(len(m.requestExternalDataResultsCounter)) <= externalID {
		return nil, 0, fmt.Errorf("externalID is out of range")
	}

	if int64(len(m.requestExternalDataResultsCounter[externalID])) <= validatorIndex {
		return nil, 0, fmt.Errorf("validatorIndex is out of range")
	}

	m.requestExternalDataResultsCounter[externalID][validatorIndex]++
	return m.externalDataResults[externalID][validatorIndex], 0, nil
}
