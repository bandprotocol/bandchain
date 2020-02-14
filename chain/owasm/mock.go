package owasm

import "fmt"

type mockExecutionEnvironment struct {
	requestID               int64
	requestedValidatorCount int64
	receivedValidatorCount  int64
	prepareBlockTime        int64
	aggregateBlockTime      int64
	validatorAddresses      [][]byte
	externalDataResults     [][][]byte
}

func (m *mockExecutionEnvironment) GetCurrentRequestID() int64 {
	return m.requestID
}

func (m *mockExecutionEnvironment) GetRequestedValidatorCount() int64 {
	return m.requestedValidatorCount
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
) ([]byte, error) {
	return m.externalDataResults[externalDataID][validatorIndex], nil
}
