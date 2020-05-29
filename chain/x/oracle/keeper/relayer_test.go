package keeper_test

import (
	"fmt"
	"testing"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmkv "github.com/tendermint/tendermint/libs/kv"
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

	capName := ibctypes.ChannelCapabilityPath(TestPortA, TestChannelA)

	testCases := []struct {
		description      string
		exec             func(log.Logger) (*app.BandApp, sdk.Context)
		expectedPass     bool
		expectedErrorLog string
	}{
		{
			"success send IBC packet",
			func(logger log.Logger) (*app.BandApp, sdk.Context) {

				chainA, chainB := createTestChains(logger)
				ctx := getContext(chainA)

				cap, err := chainA.ScopedIBCKeeper.NewCapability(ctx, capName)
				require.NoError(t, err)
				err = chainA.OracleKeeper.ScopedKeeper.ClaimCapability(ctx, cap, capName)
				require.NoError(t, err)

				createTestClient(chainA, chainB)
				createTestChainConnection(chainA, chainB)
				createTestChannel(chainA, chainB)

				ctx = getContext(chainA)
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, 1)

				return chainA, ctx
			},
			true,
			"",
		},
		{
			"failed to get channel",
			func(logger log.Logger) (*app.BandApp, sdk.Context) {
				chainA, chainB := createTestChains(logger)
				createTestChainConnection(chainA, chainB)
				ctx := getContext(chainA)
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, 1)
				return chainA, ctx
			},
			false,
			fmt.Sprintf("failed to get channel with id: %s, port: %s", TestChannelA, TestPortA),
		},
		{
			"failed to get next sequence for channel id",
			func(logger log.Logger) (*app.BandApp, sdk.Context) {
				chainA, chainB := createTestChains(logger)
				createTestChainConnection(chainA, chainB)
				createTestChannel(chainA, chainB)
				ctx := getContext(chainA)
				return chainA, ctx
			},
			false,
			fmt.Sprintf("failed to get next sequence for channel id: %s, port: %s", TestChannelA, TestPortA),
		},
		{
			"failed to get channel capability for id",
			func(logger log.Logger) (*app.BandApp, sdk.Context) {
				chainA, chainB := createTestChains(logger)
				createTestChainConnection(chainA, chainB)
				createTestChannel(chainA, chainB)

				ctx := getContext(chainA)
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, 1)
				return chainA, ctx
			},
			false,
			fmt.Sprintf("failed to get channel capability for id: %s, port: %s", TestChannelA, TestPortA),
		},
	}

	for _, testcase := range testCases {
		logger := &TestErrLogger{}
		chainA, ctx := testcase.exec(logger)

		res := types.OracleResponsePacketData{
			ClientID:      "alice",
			RequestID:     3,
			AnsCount:      1,
			RequestTime:   1589535020,
			ResolveTime:   1589535022,
			ResolveStatus: 1,
			Result:        []byte("4bb10e0000000000"),
		}

		chainA.OracleKeeper.SendOracleResponse(ctx, TestPortA, TestChannelA, res)
		if !testcase.expectedPass {
			require.Equal(t, testcase.expectedErrorLog, logger.errLog)
		} else {
			require.Equal(t, "", logger.errLog)
			events := ctx.EventManager().ABCIEvents()

			expectedEvent := sdk.Event{
				Type:       channeltypes.EventTypeSendPacket,
				Attributes: []tmkv.Pair{},
			}

			events = ctx.EventManager().ABCIEvents()
			require.Equal(t, 1, len(events))
			require.Equal(t, abci.Event(expectedEvent), events[0])
		}
	}
}
