package owasm

import (
	"fmt"
)

type mockExecEnv struct {
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

func (m *mockExecEnv) GetRequestedValidatorCount() int64 {
	return m.requestedValidatorCount
}

func (m *mockExecEnv) GetSufficientValidatorCount() int64 {
	return m.sufficientValidatorCount
}

func (m *mockExecEnv) GetReceivedValidatorCount() int64 {
	return m.receivedValidatorCount
}

func (m *mockExecEnv) GetPrepareBlockTime() int64 {
	return m.prepareBlockTime
}

func (m *mockExecEnv) GetAggregateBlockTime() int64 {
	return m.aggregateBlockTime
}

func (m *mockExecEnv) GetValidatorAddress(validatorIndex int64) ([]byte, error) {
	return m.validatorAddresses[validatorIndex], nil
}

func (m *mockExecEnv) GetMaximumResultSize() int64 {
	return m.maximumResultSize
}

func (m *mockExecEnv) GetMaximumCalldataOfDataSourceSize() int64 {
	return m.maximumCalldataOfDataSourceSize
}

func (m *mockExecEnv) RequestExternalData(
	dataSourceID int64,
	externalID int64,
	calldata []byte,
) error {
	// TODO: Figure out how to mock this elegantly.
	fmt.Printf("RequestExternalData: DataSourceID = %d, ExternalID = %d\n", dataSourceID, externalID)
	return nil
}

func (m *mockExecEnv) GetExternalData(
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
