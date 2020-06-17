package api

type RawRequest struct {
	ExternalID   int64
	DataSourceID int64
	Calldata     []byte
}

func NewRawRequest(eid int64, did int64, calldata []byte) RawRequest {
	return RawRequest{
		ExternalID:   eid,
		DataSourceID: did,
		Calldata:     calldata,
	}
}

type RawReport struct {
	ExternalID int64
	ExitCode   uint32
	Data       []byte
}

type MockEnv struct {
	Calldata    []byte
	Retdata     []byte
	rawRequests []RawRequest
}

func NewMockEnv(calldata []byte) *MockEnv {
	return &MockEnv{
		Calldata:    calldata,
		Retdata:     []byte{},
		rawRequests: []RawRequest{},
	}
}

func (env *MockEnv) GetCalldata() []byte {
	return env.Calldata
}

func (env *MockEnv) SetReturnData(data []byte) {
	env.Retdata = data
}

func (env *MockEnv) AskExternalData(eid int64, did int64, data []byte) {
	env.rawRequests = append(env.rawRequests, NewRawRequest(
		eid, did, data,
	))
}

func (env *MockEnv) GetExternalDataFull(eid int64, valIdx int64) ([]byte, int64) {
	return []byte("BEEB"), 0
}

func (env *MockEnv) GetExternalDataStatus(eid int64, vid int64) int64 {
	_, status := env.GetExternalDataFull(eid, vid)
	return status
}

func (env *MockEnv) GetExternalData(eid int64, vid int64) []byte {
	data, _ := env.GetExternalDataFull(eid, vid)
	return data
}

func (env *MockEnv) GetAskCount() int64 {
	return 0
}

func (env *MockEnv) GetMinCount() int64 {
	return 0
}

func (env *MockEnv) GetAnsCount() int64 {
	return 0
}
