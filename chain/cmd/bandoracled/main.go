package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/bandprotocol/bandchain/chain/app"
	"github.com/bandprotocol/bandchain/chain/bandlib"
	"github.com/bandprotocol/bandchain/chain/byteexec"
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

const (
	flagMaxQueryDuration = "max-query-duration"
	flagPrivKey          = "priv-key"
	flagSandboxMode      = "sandbox"
	flagExecuteEndPoint  = "execute-endpoint"
	flagChainID          = "chain-id"
)

var (
	bandClient bandlib.BandStatefulClient
	gasFlagVar = flags.GasSetting{Gas: flags.DefaultGasLimit}
	logger     log.Logger
	chainID    string
)

func getLatestRequestID() (oracle.RequestID, error) {
	cliCtx := bandClient.GetContext()
	res, _, err := cliCtx.Query("custom/oracle/request_number")
	if err != nil {
		return 0, err
	}
	var requestID oracle.RequestID
	err = cliCtx.Codec.UnmarshalJSON(res, &requestID)
	if err != nil {
		return 0, err
	}
	return requestID, nil
}

func main() {
	config := sdk.GetConfig()
	app.SetBech32AddressPrefixesAndBip44CoinType(config)
	config.Seal()
	cmd := &cobra.Command{
		Use:   "bandoracled",
		Short: "Band oracle Daemon",
		Long: strings.TrimSpace(
			`Band oracle to listen new requests from chain and reports data for execution.
Example:
$ bandoracled --node tcp://localhost:26657 --priv-key 06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b -d 60
`,
		),
		Args: cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			var err error
			logger = log.NewTMLogger(os.Stdout)
			privB, err := hex.DecodeString(viper.GetString(flagPrivKey))
			if err != nil {
				logger.Error(fmt.Sprintf("%v", err))
				return err
			}
			var priv secp256k1.PrivKeySecp256k1
			copy(priv[:], privB)

			viper.GetString(flags.FlagChainID)

			bandClient, err = bandlib.NewBandStatefulClient(
				viper.GetString(flags.FlagNode), priv, 100, 10, "Bandoracled reports", chainID,
			)
			if err != nil {
				logger.Error(fmt.Sprintf("%v", err))
				return err
			}

			currentRequestID, err := getLatestRequestID()
			if err != nil {
				logger.Error(fmt.Sprintf("%v", err))
				return err
			}

			logger.Info(fmt.Sprintf("Started at request #%d", currentRequestID))

			// Setup poll loop
			for {
				newRequestID, err := getLatestRequestID()
				if err != nil {
					logger.Error(fmt.Sprintf("%v", err))
				}
				for currentRequestID < newRequestID {
					currentRequestID++
					go handleRequest(currentRequestID)
				}
				time.Sleep(1 * time.Second)
			}
		},
	}

	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	cmd.Flags().IntP(flagMaxQueryDuration, "d", 60, "Max duration to query data")
	viper.BindPFlag(flagMaxQueryDuration, cmd.Flags().Lookup(flagMaxQueryDuration))
	cmd.Flags().String(flags.FlagFees, "", "Fees to pay along with transaction; eg: 10uband")
	viper.BindPFlag(flags.FlagFees, cmd.Flags().Lookup(flags.FlagFees))
	cmd.Flags().String(flags.FlagGasPrices, "", "Gas prices to determine the transaction fee (e.g. 10uband)")
	viper.BindPFlag(flags.FlagGasPrices, cmd.Flags().Lookup(flags.FlagGasPrices))

	// --gas can accept integers and "simulate"
	cmd.Flags().Var(&gasFlagVar, "gas", fmt.Sprintf(
		"gas limit to set per-transaction; set to %q to calculate required gas automatically (default %d)",
		flags.GasFlagAuto, flags.DefaultGasLimit,
	))
	cmd.Flags().BoolP(flagSandboxMode, "s", false, "Enable sandbox mode")
	viper.BindPFlag(flagSandboxMode, cmd.Flags().Lookup(flagSandboxMode))
	cmd.Flags().String(
		flagPrivKey,
		"06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b",
		"Private key of validator to send report transaction",
	)
	viper.BindPFlag(flagPrivKey, cmd.Flags().Lookup(flagPrivKey))

	cmd.Flags().String(
		flagExecuteEndPoint,
		"",
		"The URL of execution end-point which will receive 3 parameters (executable, timeout, calldata)",
	)
	viper.BindPFlag(flagExecuteEndPoint, cmd.Flags().Lookup(flagExecuteEndPoint))

	cmd.Flags().String(flagChainID, "", "ID of the chain to use")
	viper.BindPFlag(flagChainID, cmd.Flags().Lookup(flagChainID))

	err := cmd.Execute()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed executing CLI command: %s, exiting...", err))
		os.Exit(1)
	}
}

func handleRequest(requestID oracle.RequestID) {
	cliCtx := bandClient.GetContext()
	res, _, err := cliCtx.Query(fmt.Sprintf("custom/oracle/request/%d", requestID))
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot get request #%d. Error: %v", requestID, err))
		return
	}
	var request oracle.RequestQuerierInfo
	err = cliCtx.Codec.UnmarshalJSON(res, &request)
	if err != nil {
		logger.Error(fmt.Sprintf("Report fail on request #%d. Error: %v", requestID, err))
		return
	}

	found := false
	for _, item := range request.Request.RequestedValidators {
		if item.Equals(sdk.ValAddress(bandClient.Sender())) {
			found = true
		}
	}
	if !found {
		return
	}

	type queryParallelInfo struct {
		externalID oracle.ExternalID
		answer     []byte
		err        error
	}

	chanQueryParallelInfo := make(chan queryParallelInfo, len(request.RawDataRequests))
	for _, rawRequest := range request.RawDataRequests {
		go func(externalID oracle.ExternalID, dataSourceID oracle.DataSourceID, calldata []byte) {
			info := queryParallelInfo{externalID: externalID, answer: []byte{}, err: nil}
			res, _, err := cliCtx.Query(
				fmt.Sprintf("custom/oracle/%s/%d", oracle.QueryDataSourceByID, dataSourceID),
			)

			if err != nil {
				info.err = fmt.Errorf(
					"Cannot get data source id [%d], error: %v", dataSourceID, err,
				)
				chanQueryParallelInfo <- info
				return
			}

			var dataSource oracle.DataSourceQuerierInfo
			err = cliCtx.Codec.UnmarshalJSON(res, &dataSource)
			if err != nil {
				info.err = err
				chanQueryParallelInfo <- info
				return
			}

			var result []byte
			if viper.IsSet(flagExecuteEndPoint) {
				result, err = byteexec.RunOnAWSLambda(
					dataSource.Executable,
					time.Duration(viper.GetInt(flagMaxQueryDuration))*time.Second,
					string(calldata),
					viper.GetString(flagExecuteEndPoint),
				)
			} else {
				result, err = byteexec.RunOnDocker(
					dataSource.Executable,
					viper.IsSet(flagSandboxMode),
					time.Duration(viper.GetInt(flagMaxQueryDuration))*time.Second,
					string(calldata),
				)
			}

			if err != nil {
				info.err = fmt.Errorf(
					"Execute error on data source id [%d], error: %v", dataSourceID, err,
				)
				chanQueryParallelInfo <- info
				return
			}

			info.answer = []byte(strings.TrimSpace(string(result)))
			chanQueryParallelInfo <- info
		}(rawRequest.ExternalID,
			rawRequest.RawDataRequest.DataSourceID,
			rawRequest.RawDataRequest.Calldata,
		)
	}

	reports := make([]oracle.RawDataReportWithID, 0)
	for i := 0; i < len(request.RawDataRequests); i++ {
		info := <-chanQueryParallelInfo
		if info.err != nil {
			logger.Error(fmt.Sprintf("Report fail on request #%d. Error: %v", requestID, info.err))
			return
		}
		reports = append(reports, oracle.NewRawDataReportWithID(info.externalID, 0, info.answer))
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ExternalDataID < reports[j].ExternalDataID
	})

	if err != nil {
		logger.Error(fmt.Sprintf("Send report fail on request #%d. Error: %v", requestID, err))
		return
	}
	tx, err := bandClient.SendTransaction(
		oracle.NewMsgReportData(requestID, reports, sdk.ValAddress(bandClient.Sender()), bandClient.Sender()),
		gasFlagVar.Gas, viper.GetString(flags.FlagFees), viper.GetString(flags.FlagGasPrices),
	)

	if err != nil {
		logger.Error(fmt.Sprintf("Report fail on request #%d. Error: %v", requestID, err))
		return
	}
	logger.Info(fmt.Sprintf("Report on request #%d successfully. Tx: %v", requestID, tx))
}
