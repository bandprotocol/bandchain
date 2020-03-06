package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/bandlib"
	"github.com/bandprotocol/d3n/chain/byteexec"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

const (
	flagMaxQueryDuration = "max-query-duration"
	flagPrivKey          = "priv-key"
)

var (
	bandClient bandlib.BandStatefulClient
)

func getLatestRequestID() (int64, error) {
	cliCtx := bandClient.GetContext()
	res, _, err := cliCtx.Query("custom/zoracle/request_number")
	if err != nil {
		return 0, err
	}
	var requestID int64
	err = cliCtx.Codec.UnmarshalJSON(res, &requestID)
	if err != nil {
		return 0, err
	}
	return requestID, nil
}

func main() {
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

			privB, err := hex.DecodeString(viper.GetString(flagPrivKey))
			if err != nil {
				return err
			}
			var priv secp256k1.PrivKeySecp256k1
			copy(priv[:], privB)

			bandClient, err = bandlib.NewBandStatefulClient(
				viper.GetString(flags.FlagNode),
				priv,
			)
			if err != nil {
				return err
			}

			currentRequestID, err := getLatestRequestID()
			if err != nil {
				return err
			}

			// Setup poll loop
			for {
				newRequestID, err := getLatestRequestID()
				if err != nil {
					log.Println("Cannot get request number error: ", err.Error())
				}
				for currentRequestID < newRequestID {
					currentRequestID++
					go handleRequestAndLog(currentRequestID)
				}
				time.Sleep(1 * time.Second)
			}
		},
	}

	cmd.Flags().String(flags.FlagNode, "tcp://localhost:26657", "<host>:<port> to Tendermint RPC interface for this chain")
	viper.BindPFlag(flags.FlagNode, cmd.Flags().Lookup(flags.FlagNode))
	cmd.Flags().IntP(flagMaxQueryDuration, "d", 60, "Max duration to query data")
	viper.BindPFlag(flagMaxQueryDuration, cmd.Flags().Lookup(flagMaxQueryDuration))
	cmd.Flags().String(
		flagPrivKey,
		"06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b",
		"Private key of validator to send report transaction",
	)
	viper.BindPFlag(flagPrivKey, cmd.Flags().Lookup(flagPrivKey))
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func handleRequest(requestID int64) (sdk.TxResponse, error) {
	cliCtx := bandClient.GetContext()
	res, _, err := cliCtx.Query(fmt.Sprintf("custom/zoracle/request/%d", requestID))
	if err != nil {
		return sdk.TxResponse{}, err
	}
	var request zoracle.RequestQuerierInfo
	err = cliCtx.Codec.UnmarshalJSON(res, &request)
	if err != nil {
		return sdk.TxResponse{}, err
	}

	type queryParallelInfo struct {
		externalID int64
		answer     []byte
		err        error
	}

	chanQueryParallelInfo := make(chan queryParallelInfo, len(request.RawDataRequests))
	for _, rawRequest := range request.RawDataRequests {
		go func(externalID, dataSourceID int64, calldata []byte) {
			info := queryParallelInfo{externalID: externalID, answer: []byte{}, err: nil}
			res, _, err := cliCtx.Query(
				fmt.Sprintf("custom/zoracle/%s/%d", zoracle.QueryDataSourceByID, dataSourceID),
			)

			if err != nil {
				info.err = fmt.Errorf(
					"handleRequest: Cannot get script id [%d], error: %v", dataSourceID, err,
				)
				chanQueryParallelInfo <- info
				return
			}

			var dataSource zoracle.DataSourceQuerierInfo
			err = cliCtx.Codec.UnmarshalJSON(res, &dataSource)
			if err != nil {
				info.err = err
				chanQueryParallelInfo <- info
				return
			}

			result, err := byteexec.RunOnDocker(
				dataSource.Executable,
				time.Duration(viper.GetInt(flagMaxQueryDuration))*time.Second,
				string(calldata),
			)
			if err != nil {
				info.err = fmt.Errorf(
					"handleRequest: Execute error on data request id [%d], error: %v", dataSourceID, err,
				)
				chanQueryParallelInfo <- info
				return
			}

			info.answer = []byte(strings.TrimSpace(string(result)))
			chanQueryParallelInfo <- info
		}(rawRequest.ExternalID, rawRequest.RawDataRequest.DataSourceID, rawRequest.RawDataRequest.Calldata)
	}

	reports := make([]zoracle.RawDataReport, 0)
	for i := 0; i < len(request.RawDataRequests); i++ {
		info := <-chanQueryParallelInfo
		if info.err != nil {
			return sdk.TxResponse{}, info.err
		}
		reports = append(reports, zoracle.NewRawDataReport(info.externalID, info.answer))
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ExternalDataID < reports[j].ExternalDataID
	})

	return bandClient.SendTransaction(
		zoracle.NewMsgReportData(requestID, reports, sdk.ValAddress(bandClient.Sender())),
		1000000, "", "", "",
		flags.BroadcastSync,
	)
}

func handleRequestAndLog(requestID int64) {
	fmt.Println(handleRequest(requestID))
}
