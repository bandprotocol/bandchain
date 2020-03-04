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
	"github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/bandprotocol/d3n/chain/bandlib"
	"github.com/bandprotocol/d3n/chain/byteexec"
	"github.com/bandprotocol/d3n/chain/x/zoracle"
)

const (
	limitTimeOut = 1 * time.Minute
)

var (
	bandClient bandlib.BandStatefulClient
	nodeURI    = getEnv("NODE_URI", "http://localhost:26657")
	privS      = getEnv("PRIVATE_KEY", "06be35b56b048c5a6810a47e2ef612eaed735ccb0d7ea4fc409f23f1d1a16e0b")
)

func getEnv(key, defaultValue string) string {
	tmp := os.Getenv(key)
	if tmp == "" {
		return defaultValue
	}
	return tmp
}

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
	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)

	var err error
	bandClient, err = bandlib.NewBandStatefulClient(nodeURI, priv)
	if err != nil {
		panic(err)
	}

	currentRequestID, err := getLatestRequestID()
	if err != nil {
		panic(err)
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
		go func(externalID , dataSourceID int64, calldata []byte) {
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

			result, err := byteexec.RunOnDocker(dataSource.Executable, limitTimeOut, string(calldata))
			if err != nil {
				info.err = fmt.Errorf(
					"handleRequest: Execute error on data request id [%d], error: %v", dataSourceID, err,
				)
				chanQueryParallelInfo <- info
				return
			}

			info.answer = []byte(strings.TrimSpace(string(result)))
			chanQueryParallelInfo <- info
		}(int64(rawRequest.ExternalID), int64(rawRequest.RawDataRequest.DataSourceID), rawRequest.RawDataRequest.Calldata)
	}

	reports := make([]zoracle.RawDataReport, 0)
	for i := 0; i < len(request.RawDataRequests); i++ {
		info := <-chanQueryParallelInfo
		if info.err != nil {
			return sdk.TxResponse{}, info.err
		}
		reports = append(reports, zoracle.NewRawDataReport(zoracle.ExternalID(info.externalID), info.answer))
	}

	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ExternalDataID < reports[j].ExternalDataID
	})

	return bandClient.SendTransaction(
		zoracle.NewMsgReportData(zoracle.RequestID(requestID), reports, sdk.ValAddress(bandClient.Sender())),
		1000000, "", "", "",
		flags.BroadcastSync,
	)
}

func handleRequestAndLog(requestID int64) {
	fmt.Println(handleRequest(requestID))
}
