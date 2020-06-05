package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func creatDeafaultServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		ret := externalExecutionResponse{
			Returncode: 0,
			Stdout:     "BEEB",
			Stderr:     "Stderr",
		}
		json.NewEncoder(res).Encode(ret)
	}))
}

func createResponseNotOkSenarioServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
	}))
}

func createCannotDecodeJsonSenarioServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("invalid bytes"))
	}))
}

func creatExecuteFailSenarioServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		ret := externalExecutionResponse{
			Returncode: 1,
			Stdout:     "BEEB",
			Stderr:     "Stderr",
		}
		json.NewEncoder(res).Encode(ret)
	}))
}

func getLog() *Logger {
	logLevel, _ := log.AllowLevel("debug")
	return NewLogger(logLevel)
}

func TestExecuteSuccess(t *testing.T) {
	testServer := creatDeafaultServer()
	defer func() { testServer.Close() }()

	executor := &lambdaExecutor{URL: testServer.URL}
	res, exitcode := executor.Execute(getLog(), []byte("executable"), 1*time.Second, "calldata")

	require.Equal(t, uint32(0), exitcode)
	require.Equal(t, []byte("BEEB"), res)
}

func TestExecuteBadUrlFail(t *testing.T) {
	testServer := creatDeafaultServer()
	defer func() { testServer.Close() }()

	executor := &lambdaExecutor{URL: "www.beeb.com"} // bad url
	res, exitcode := executor.Execute(getLog(), []byte("executable"), 1*time.Second, "calldata")

	require.Equal(t, uint32(255), exitcode)
	require.Equal(t, []byte("EXECUTION_ERROR"), res)
}

func TestExecuteDecodeStructFail(t *testing.T) {
	testServer := createCannotDecodeJsonSenarioServer()
	defer func() { testServer.Close() }()

	executor := &lambdaExecutor{URL: testServer.URL}
	res, exitcode := executor.Execute(getLog(), []byte("executable"), 1*time.Second, "calldata")
	require.Equal(t, uint32(255), exitcode)
	require.Equal(t, []byte("EXECUTION_ERROR"), res)
}

func TestExecuteResponseNotOk(t *testing.T) {
	testServer := createResponseNotOkSenarioServer()
	defer func() { testServer.Close() }()

	executor := &lambdaExecutor{URL: testServer.URL}
	res, exitcode := executor.Execute(getLog(), []byte("executable"), 1*time.Second, "calldata")
	require.Equal(t, uint32(255), exitcode)
	require.Equal(t, []byte("EXECUTION_ERROR"), res)
}

func TestExecuteFail(t *testing.T) {
	testServer := creatExecuteFailSenarioServer()
	defer func() { testServer.Close() }()

	executor := &lambdaExecutor{URL: testServer.URL}
	res, exitcode := executor.Execute(getLog(), []byte("executable"), 1*time.Second, "calldata")
	require.Equal(t, uint32(1), exitcode)
	require.Equal(t, []byte{}, res)
}
