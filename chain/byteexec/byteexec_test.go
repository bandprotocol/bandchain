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
	require.Equal(t, []byte("hello world\n"), result)
}
