package byteexec

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExecByte(t *testing.T) {
	// Simple bash script
	execByte, _ := hex.DecodeString("23212f62696e2f626173680a0a6563686f2024310a")
	result, err := RunOnLocal(execByte, 10*time.Second, "hello world")
	require.Nil(t, err)
	require.Equal(t, []byte("hello\n"), result)
}

func TestExecByte2(t *testing.T) {
	// Simple bash script
	execByte, _ := hex.DecodeString("23212f62696e2f62617368aa6563686f202432a")
	result, err := RunOnLocal(execByte, 10*time.Second, "aS world")
	require.Nil(t, err)
	require.Equal(t, []byte("world\n"), result)
}

// 23212f62696e2f626173680a0a6563686f2024322b24310a
func TestExecByte3(t *testing.T) {
	// Simple bash script
	execByte, _ := hex.DecodeString("23212f62696e2f626173680a0a6563686f2024322b24310a")
	result, err := RunOnLocal(execByte, 10*time.Second, "HELLO WORLD")
	require.Nil(t, err)
	require.Equal(t, []byte("HELLOWORLD\n"), result)
}
