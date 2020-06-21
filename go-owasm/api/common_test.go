package api

import (
	"io/ioutil"
	"os"
	"os/exec"
)

// wat2wasm compiles the given Wat content to Wasm, relying on the host's wat2wasm program.
func wat2wasm(wat []byte) []byte {
	inputFile, err := ioutil.TempFile("", "input")
	if err != nil {
		panic(err)
	}
	defer os.Remove(inputFile.Name())
	outputFile, err := ioutil.TempFile("", "output")
	if err != nil {
		panic(err)
	}
	defer os.Remove(outputFile.Name())
	if _, err := inputFile.Write(wat); err != nil {
		panic(err)
	}
	if err := exec.Command("wat2wasm", inputFile.Name(), "-o", outputFile.Name()).Run(); err != nil {
		panic(err)
	}
	output, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		panic(err)
	}
	return output
}

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

func (env *MockEnv) SetReturnData(data []byte) error {
	env.Retdata = data
	return nil
}

func (env *MockEnv) AskExternalData(eid int64, did int64, data []byte) error {
	env.rawRequests = append(env.rawRequests, NewRawRequest(
		eid, did, data,
	))
	return nil
}

func (env *MockEnv) GetExternalDataFull(eid int64, valIdx int64) ([]byte, int64) {
	return []byte("BEEB"), 0
}

func (env *MockEnv) GetExternalDataStatus(eid int64, vid int64) (int64, error) {
	_, status := env.GetExternalDataFull(eid, vid)
	return status, nil
}

func (env *MockEnv) GetExternalData(eid int64, vid int64) ([]byte, error) {
	data, _ := env.GetExternalDataFull(eid, vid)
	return data, nil
}

func (env *MockEnv) GetAskCount() int64 {
	return 0
}

func (env *MockEnv) GetMinCount() int64 {
	return 0
}

func (env *MockEnv) GetAnsCount() (int64, error) {
	return 0, nil
}
