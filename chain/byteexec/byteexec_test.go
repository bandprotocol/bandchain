package byteexec

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecByte(t *testing.T) {
	// Simple bash script
	execByte, _ := hex.DecodeString("23212f62696e2f626173680a0a6563686f2024310a")
	ex, err := New(execByte)
	require.Nil(t, err)

	defer ex.Close()
	cmd := ex.Command("hello world")
	result, err := cmd.Output()
	require.Nil(t, err)
	require.Equal(t, []byte("hello world\n"), result)
}
