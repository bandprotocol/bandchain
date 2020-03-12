package byteexec

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecByte(t *testing.T) {
	// Simple bash script
	/*
		!/bin/bash

		echo $1

	*/
	execByte, _ := hex.DecodeString("23212f62696e2f626173680a0a6563686f2024310a")
	result, err := RunOnLocal(execByte, 10*time.Second, "hello world")
	require.Nil(t, err)
	require.Equal(t, []byte("hello\n"), result)
}

func TestExecByte2(t *testing.T) {
	// Simple bash script
	/*
		!/bin/bash

		echo $2

	*/
	execByte, _ := hex.DecodeString("23212f62696e2f626173680a0a6563686f2024320a")
	result, err := RunOnLocal(execByte, 10*time.Second, "aS world")
	require.Nil(t, err)
	require.Equal(t, []byte("world\n"), result)
}

// 23212f62696e2f626173680a0a6563686f2024322b24310a
func TestExecByte3(t *testing.T) {
	// batch script that can sum two number
	/*
		#!/bin/bash
		num1=$1
		num2=$2
		echo $((num1+num2))
	*/
	execByte, _ := hex.DecodeString("23212f62696e2f626173680a6e756d313d24310a6e756d323d24320a6563686f202428286e756d312b6e756d3229290a")
	result, err := RunOnLocal(execByte, 10*time.Second, "4 14")
	require.Nil(t, err)
	require.Equal(t, []byte("18\n"), result)
}
