package owasm

type mockExecutionEnvironment struct {
}

func (m *mockExecutionEnvironment) GetCurrentRequestID() int64 {
	return 0
}

func (m *mockExecutionEnvironment) GetRequestedValidatorCount() int64 {
	return 0
}

func (m *mockExecutionEnvironment) GetReceivedValidatorCount() int64 {
	return 0
}

func (m *mockExecutionEnvironment) GetPrepareBlockTime() int64 {
	return 0
}

func (m *mockExecutionEnvironment) GetAggregateBlockTime() int64 {
	return 0
}

func (m *mockExecutionEnvironment) GetValidatorPubKey(validatorIndex int64) ([]byte, error) {
	return []byte{}, nil
}

func (m *mockExecutionEnvironment) RequestExternalData(
	dataSourceID int64,
	externalDataID int64,
	calldata []byte,
) error {
	return nil
}

func (m *mockExecutionEnvironment) GetExternalData(
	externalDataID int64,
	validatorIndex int64,
) ([]byte, error) {
	return []byte{}, nil
}

func NewMockExecutionEnvironment() ExecutionEnvironment {
	return &mockExecutionEnvironment{}
}
