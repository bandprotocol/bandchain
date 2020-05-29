package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/bandprotocol/bandchain/chain/app"
	bandapp "github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/simapp"
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	connectiontypes "github.com/cosmos/cosmos-sdk/x/ibc/03-connection/types"
	channelexported "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/exported"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctmtypes "github.com/cosmos/cosmos-sdk/x/ibc/07-tendermint/types"
	commitmenttypes "github.com/cosmos/cosmos-sdk/x/ibc/23-commitment/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
)

const (
	ChainIDA                                            = "chainA"
	ChainIDB                                            = "chainB"
	TestClientIDA                                       = "clientA"
	TestClientIDB                                       = "clientB"
	TestPortA                                           = "testporta"
	TestPortB                                           = "testportb"
	TestChannelA                                        = "testchannela"
	TestChannelB                                        = "testchannelb"
	TestConnectionA                                     = "connectionAtoB"
	TestConnectionB                                     = "connectionBtoA"
	TrustingPeriod                        time.Duration = time.Hour * 24 * 7 * 2
	UbdPeriod                             time.Duration = time.Hour * 24 * 7 * 3
	MaxClockDrift                         time.Duration = time.Second * 10
	DefaultPacketTimeoutHeight                          = 0
	DefaultPacketTimeoutTimestampDuration               = uint64(600 * time.Second)
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

func newDefaultRequest() types.Request {
	return types.NewRequest(
		1,
		[]byte("calldata"),
		[]sdk.ValAddress{Validator1.ValAddress, Validator2.ValAddress},
		2,
		0,
		1581503227,
		"clientID",
		nil,
		[]types.ExternalID{42},
	)
}

func createTestChains(logger log.Logger) (*bandapp.BandApp, *bandapp.BandApp) {
	appA := simapp.NewSimApp(ChainIDA, logger)
	appB := simapp.NewSimApp(ChainIDB, logger)
	return appA, appB
}

func getContext(chain *bandapp.BandApp) sdk.Context {
	now := time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC)

	privVal := tmtypes.NewMockPV()
	signers := []tmtypes.PrivValidator{privVal}
	pubKey, err := privVal.GetPubKey()
	if err != nil {
		panic(err)
	}
	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	header := ibctmtypes.CreateTestHeader(chain.Name(), 1, now, valSet, signers)

	return chain.NewContext(false, abci.Header{
		ChainID: header.ChainID,
		Height:  header.Height,
		Time:    now,
	})
}

func createTestClient(chainA *bandapp.BandApp, chainB *bandapp.BandApp) error {
	oldCtx := getContext(chainB)

	// Commit and create a new block on client to get a fresh CommitID
	chainB.Commit()
	commitID := chainB.LastCommitID()

	chainB.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: oldCtx.BlockHeight() + 1, Time: oldCtx.BlockTime().Add(time.Minute)}})

	// Set HistoricalInfo on client chain after Commit
	newCtxClient := getContext(chainB)

	// Prepare validator and signers for client chain
	privVal := tmtypes.NewMockPV()

	pubKey, err := privVal.GetPubKey()
	if err != nil {
		panic(err)
	}

	validator := tmtypes.NewValidator(pubKey, 1)
	valSet := tmtypes.NewValidatorSet([]*tmtypes.Validator{validator})
	signers := []tmtypes.PrivValidator{privVal}

	stakingValidator := staking.NewValidator(
		sdk.ValAddress(valSet.Validators[0].Address), valSet.Validators[0].PubKey, staking.Description{},
	)
	stakingValidator.Status = sdk.Bonded
	stakingValidator.Tokens = sdk.NewInt(1000000)
	stakingValidators := []staking.Validator{stakingValidator}

	histInfo := stakingtypes.HistoricalInfo{
		Header: abci.Header{
			AppHash: commitID.Hash,
		},
		Valset: stakingValidators,
	}
	chainB.StakingKeeper.SetHistoricalInfo(newCtxClient, newCtxClient.BlockHeader().Height, histInfo)

	// Create target context
	ctxTarget := getContext(chainA)

	// Create client
	header := ibctmtypes.CreateTestHeader(ChainIDB, newCtxClient.BlockHeader().Height+1, newCtxClient.BlockTime().Add(time.Minute), valSet, signers)
	clientState, err := ibctmtypes.Initialize(TestClientIDB, TrustingPeriod, UbdPeriod, MaxClockDrift, header)
	if err != nil {
		return err
	}

	_, err = chainA.IBCKeeper.ClientKeeper.CreateClient(ctxTarget, clientState, header.ConsensusState())
	if err != nil {
		return err
	}
	return nil
}

func createTestChainConnection(chainA *bandapp.BandApp, chainB *bandapp.BandApp) {
	counterParty := connectiontypes.NewCounterparty(TestConnectionA, TestConnectionB, commitmenttypes.NewMerklePrefix(chainA.IBCKeeper.ConnectionKeeper.GetCommitmentPrefix().Bytes()))
	conn := connectiontypes.NewConnectionEnd(3, TestClientIDB, counterParty, connectiontypes.GetCompatibleVersions())
	ctx := chainA.NewContext(false, abci.Header{})
	chainA.IBCKeeper.ConnectionKeeper.SetConnection(ctx, TestConnectionA, conn)
}

func createTestChannel(chainA *bandapp.BandApp, chainB *bandapp.BandApp) {
	counterpart := channeltypes.NewCounterparty(TestPortB, TestChannelB)
	channel := channeltypes.NewChannel(channelexported.OPEN, channelexported.ORDERED, counterpart, []string{TestConnectionA}, "1.0")
	ctx := chainA.NewContext(false, abci.Header{})
	chainA.IBCKeeper.ChannelKeeper.SetChannel(ctx, TestPortA, TestChannelA, channel)
}

func TestSendOracleResponse(t *testing.T) {

	capName := ibctypes.ChannelCapabilityPath(TestPortA, TestChannelA)
	seq := uint64(1)

	testCases := []struct {
		description      string
		exec             func(log.Logger) (*app.BandApp, sdk.Context)
		expectedPass     bool
		expectedErrorLog string
	}{
		{
			"successfully send IBC packet",
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
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, seq)

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
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, seq)
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
				chainA.IBCKeeper.ChannelKeeper.SetNextSequenceSend(ctx, TestPortA, TestChannelA, seq)
				return chainA, ctx
			},
			false,
			fmt.Sprintf("failed to get channel capability for id: %s, port: %s", TestChannelA, TestPortA),
		},
	}

	for _, testcase := range testCases {
		logger := &TestErrLogger{}
		chainA, ctx := testcase.exec(logger)
		ctx = ctx.WithBlockTime(time.Unix(1581588400, 0))
		ctx = ctx.WithBlockHeight(2)

		packet := types.OracleResponsePacketData{
			ClientID:      "alice",
			RequestID:     3,
			AnsCount:      1,
			RequestTime:   1589535020,
			ResolveTime:   1589535022,
			ResolveStatus: 1,
			Result:        []byte("beeb"),
		}

		chainA.OracleKeeper.SendOracleResponse(ctx, TestPortA, TestChannelA, packet)
		if !testcase.expectedPass {
			require.Equal(t, testcase.expectedErrorLog, logger.errLog)
		} else {
			require.Equal(t, "", logger.errLog)
			events := ctx.EventManager().Events()

			expectedEvents := sdk.Events{
				sdk.NewEvent(
					channeltypes.EventTypeSendPacket,
					sdk.NewAttribute(channeltypes.AttributeKeyData, fmt.Sprintf("%s", packet.GetBytes())),
					sdk.NewAttribute(channeltypes.AttributeKeyTimeoutHeight, fmt.Sprintf("%d", DefaultPacketTimeoutHeight)),
					sdk.NewAttribute(channeltypes.AttributeKeyTimeoutTimestamp, fmt.Sprintf("%d", uint64(ctx.BlockTime().UnixNano())+DefaultPacketTimeoutTimestampDuration)),
					sdk.NewAttribute(channeltypes.AttributeKeySequence, fmt.Sprintf("%d", seq)),
					sdk.NewAttribute(channeltypes.AttributeKeySrcPort, TestPortA),
					sdk.NewAttribute(channeltypes.AttributeKeySrcChannel, TestChannelA),
					sdk.NewAttribute(channeltypes.AttributeKeyDstPort, TestPortB),
					sdk.NewAttribute(channeltypes.AttributeKeyDstChannel, TestChannelB),
				),
			}

			require.Equal(t, expectedEvents, events)
		}
	}
}
