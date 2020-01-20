package main

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bandprotocol/d3n/chain/x/zoracle"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	common "github.com/tendermint/tendermint/libs/common"
)

func TestExecWithTimeoutSuccess(t *testing.T) {
	command := Command{Cmd: "ping", Arguments: []string{"-c 1", "-i 1", "8.8.8.8"}}
	_, err := execWithTimeout(command, 2000)
	require.Nil(t, err)
}

func TestExecWithTimeoutFailTimeout(t *testing.T) {
	command := Command{Cmd: "ping", Arguments: []string{"-c 2", "-i 1", "8.8.8.8"}}
	_, err := execWithTimeout(command, 1000)
	require.EqualError(t, err, "Command timed out")
}

func TestHandleRequestShouldNotTimeout(t *testing.T) {
	commands := []Command{Command{Cmd: "ping", Arguments: []string{"-c 1", "-i 1", "8.8.8.8"}}}
	value, err := json.Marshal(&commands)
	require.Nil(t, err)

	event := abci.Event{}
	event.Attributes = []common.KVPair{common.KVPair{Key: ([]byte)(zoracle.AttributeKeyPrepare), Value: ([]byte)(fmt.Sprintf("%x", value))}}
	_, err = handleRequest(&event)
	require.EqualError(t, err, "handleRequest send tx fail : ABCIQuery: Post http:: http: no Host in request URL")
}

func TestHandleRequestShouldTimeout(t *testing.T) {
	commands := []Command{Command{Cmd: "ping", Arguments: []string{"-c 4", "-i 1", "8.8.8.8"}}}
	value, err := json.Marshal(&commands)
	require.Nil(t, err)

	event := abci.Event{}
	event.Attributes = []common.KVPair{common.KVPair{Key: ([]byte)(zoracle.AttributeKeyPrepare), Value: ([]byte)(fmt.Sprintf("%x", value))}}
	_, err = handleRequest(&event)
	require.EqualError(t, err, "handleRequest query err with command: ping [-c 4 -i 1 8.8.8.8], error: Command timed out")
}
