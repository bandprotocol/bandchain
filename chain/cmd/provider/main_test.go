package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecWithTimeoutSuccess(t *testing.T) {
	command := Command{Cmd: "sleep", Arguments: []string{"0.1"}}
	_, err := execWithTimeout(command, 200*time.Millisecond)
	require.Nil(t, err)
}

func TestExecWithTimeoutFailTimeout(t *testing.T) {
	command := Command{Cmd: "sleep", Arguments: []string{"0.2"}}
	_, err := execWithTimeout(command, 100*time.Millisecond)
	require.EqualError(t, err, "Command timed out")
}
