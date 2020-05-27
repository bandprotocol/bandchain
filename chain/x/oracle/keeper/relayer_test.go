package keeper_test

import (
	"fmt"
	"testing"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

type TestErrLogger struct {
	errLog string
}

func (t *TestErrLogger) Debug(msg string, keyvals ...interface{}) {}

func (t *TestErrLogger) Info(msg string, keyvals ...interface{}) {}

func (t *TestErrLogger) Error(msg string, keyvals ...interface{}) {
	t.errLog = t.errLog + msg
}

func (t *TestErrLogger) With(keyvals ...interface{}) log.Logger {
	return t
}

func TestSendOracleResponse(t *testing.T) {

	testCases := []struct {
		description      string
		exec             func(log.Logger) (*app.BandApp, *app.BandApp, sdk.Context)
		expectedPass     bool
		expectedErrorLog string
	}{
		{
			"failed to get channel",
			func(logger log.Logger) (*app.BandApp, *app.BandApp, sdk.Context) {

				fmt.Println("start")
				chainA, chainB := createTestChains(logger)
				fmt.Println("test chain")

				createTestChainConnection(chainA, chainB)
				fmt.Println("connection")
				// createTestChannel(chainA, chainB)
				fmt.Println("channel")

				ctx := getContext(chainA)
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, 1)

				return chainA, chainB, ctx
			},
			false,
			fmt.Sprintf("failed to get channel with id: %s, port: %s", TestChannelA, TestPortA),
		},
		{
			"failed to get channel",
			func(logger log.Logger) (*app.BandApp, *app.BandApp, sdk.Context) {
				chainA, chainB := createTestChains(logger)
				createTestChainConnection(chainA, chainB)
				ctx := getContext(chainA)
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, 1)
				return chainA, chainB, ctx
			},
			false,
			fmt.Sprintf("failed to get channel with id: %s, port: %s", TestChannelA, TestPortA),
		},
		{
			"failed to get next sequence for channel id",
			func(logger log.Logger) (*app.BandApp, *app.BandApp, sdk.Context) {
				chainA, chainB := createTestChains(logger)
				createTestChainConnection(chainA, chainB)
				createTestChannel(chainA, chainB)
				ctx := getContext(chainA)
				return chainA, chainB, ctx
			},
			false,
			fmt.Sprintf("failed to get next sequence for channel id: %s, port: %s", TestChannelA, TestPortA),
		},
	}

	for _, testcase := range testCases {
		logger := &TestErrLogger{}
		chainA, _, ctx := testcase.exec(logger)

		res := types.OracleResponsePacketData{
			ClientID:      "alice",
			RequestID:     3,
			AnsCount:      1,
			RequestTime:   1589535020,
			ResolveTime:   1589535022,
			ResolveStatus: 1,
			Result:        []byte("4bb10e0000000000"),
		}

		fmt.Println("logger.errLog", logger.errLog)
		chainA.OracleKeeper.SendOracleResponse(ctx, TestPortA, TestChannelA, res)
		if !testcase.expectedPass {
			require.Equal(t, testcase.expectedErrorLog, logger.errLog)
		}
	}
}
