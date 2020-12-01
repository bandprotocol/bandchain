package proof

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncodeTime(t *testing.T) {
	require.Equal(t, hexToBytes("08d78dd9fd0510c4a1aae301"), encodeTime(parseTime("2020-11-19T10:20:07.476745924Z")))
}
