package executor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseExecutor(t *testing.T) {
	name, url, timeout, err := parseExecutor("beeb:www.beebprotocol.com?timeout=3s")
	require.Equal(t, name, "beeb")
	require.Equal(t, timeout, 3*time.Second)
	require.Equal(t, url, "www.beebprotocol.com")
	require.NoError(t, err)

	name, url, timeout, err = parseExecutor("beeb2:www.beeb.com/anna/kondanna?timeout=300ms")
	require.Equal(t, name, "beeb2")
	require.Equal(t, timeout, 300*time.Millisecond)
	require.Equal(t, url, "www.beeb.com/anna/kondanna")
	require.NoError(t, err)

	name, url, timeout, err = parseExecutor("beeb3:https://bandprotocol.com/gg/gg2/bandchain?timeout=1s300ms")
	require.Equal(t, name, "beeb3")
	require.Equal(t, timeout, 1*time.Second+300*time.Millisecond)
	require.Equal(t, url, "https://bandprotocol.com/gg/gg2/bandchain")
	require.NoError(t, err)
}

func TestParseExecutorWithoutRawQuery(t *testing.T) {
	_, _, _, err := parseExecutor("beeb:www.beebprotocol.com")
	require.EqualError(t, err, "Invalid timeout, executor requires query timeout")
}

func TestParseExecutorInvalidExecutorError(t *testing.T) {
	_, _, _, err := parseExecutor("beeb")
	require.EqualError(t, err, "Invalid executor, cannot parse executor: beeb")
}

func TestParseExecutorInvalidTimeoutError(t *testing.T) {
	_, _, _, err := parseExecutor("beeb:www.beebprotocol.com?timeout=beeb")
	require.EqualError(t, err, "Invalid timeout, cannot parse duration with error: time: invalid duration \"beeb\"")
}
