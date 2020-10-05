package common

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"

	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

func getData(cliCtx context.CLIContext, bz []byte, ptr interface{}) error {
	var result types.QueryResult
	if err := json.Unmarshal(bz, &result); err != nil {
		return err
	}
	return cliCtx.Codec.UnmarshalJSON(result.Result, ptr)
}

func queryReporters(route string, cliCtx context.CLIContext, validator sdk.ValAddress) ([]sdk.AccAddress, int64, error) {
	bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s/%s", route, types.QueryReporters, validator))
	if err != nil {
		return nil, 0, err
	}
	var reporters []sdk.AccAddress
	err = getData(cliCtx, bz, &reporters)
	if err != nil {
		return nil, 0, err
	}
	return reporters, height, nil
}

func queryParams(route string, cliCtx context.CLIContext) (types.Params, int64, error) {
	bz, height, err := cliCtx.Query(fmt.Sprintf("custom/%s/%s", route, types.QueryParams))
	if err != nil {
		return types.Params{}, 0, err
	}
	var params types.Params
	err = getData(cliCtx, bz, &params)
	if err != nil {
		return types.Params{}, 0, err
	}
	return params, height, nil
}

// TODO: Refactor this code with yoda
type VerificationMessage struct {
	ChainID    string           `json:"chain_id"`
	Validator  sdk.ValAddress   `json:"validator"`
	RequestID  types.RequestID  `json:"request_id"`
	ExternalID types.ExternalID `json:"external_id"`
}

func NewVerificationMessage(
	chainID string, validator sdk.ValAddress, requestID types.RequestID, externalID types.ExternalID,
) VerificationMessage {
	return VerificationMessage{
		ChainID:    chainID,
		Validator:  validator,
		RequestID:  requestID,
		ExternalID: externalID,
	}
}

func (msg VerificationMessage) GetSignBytes() []byte {
	return sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg))
}

type VerificationResult struct {
	ChainID      string             `json:"chain_id"`
	Validator    sdk.ValAddress     `json:"validator"`
	RequestID    types.RequestID    `json:"request_id,string"`
	ExternalID   types.ExternalID   `json:"external_id,string"`
	DataSourceID types.DataSourceID `json:"data_source_id,string"`
}

func VerifyRequest(
	route string, cliCtx context.CLIContext, chainID string, requestID types.RequestID,
	externalID types.ExternalID, validator sdk.ValAddress, reporterPubkey crypto.PubKey, signature []byte,
) ([]byte, int64, error) {
	// Verify chain id
	if cliCtx.ChainID != chainID {
		return nil, 0, fmt.Errorf("Invalid Chain ID; expect %s, got %s", cliCtx.ChainID, chainID)
	}

	// Verify signature
	if !reporterPubkey.VerifyBytes(
		NewVerificationMessage(
			chainID, validator, requestID, externalID,
		).GetSignBytes(),
		signature,
	) {
		return nil, 0, fmt.Errorf("Signature verification failed")
	}

	// Check reporters
	reporters, _, err := queryReporters(route, cliCtx, validator)
	if err != nil {
		return nil, 0, err
	}

	reporter := sdk.AccAddress(reporterPubkey.Address())

	isReporter := false
	for _, r := range reporters {
		if reporter.Equals(r) {
			isReporter = true
		}
	}
	if !isReporter {
		return nil, 0, fmt.Errorf("%s is not an authorized report of %s", reporter, validator)
	}

	request, height, err := queryRequest(route, cliCtx, requestID)
	if err != nil {
		return nil, 0, err
	}

	// Verify that this validator has been assigned to report this request
	assigned := false
	for _, requestedVal := range request.Request.RequestedValidators {
		if validator.Equals(requestedVal) {
			assigned = true
			break
		}
	}
	if !assigned {
		return nil, 0, fmt.Errorf("%s is not assigned for request ID %d", validator, requestID)
	}

	// Verify this request need this external id
	dataSourceID := types.DataSourceID(0)
	for _, rawRequest := range request.Request.RawRequests {
		if rawRequest.ExternalID == externalID {
			dataSourceID = rawRequest.DataSourceID
			break
		}
	}
	if dataSourceID == types.DataSourceID(0) {
		return nil, 0, fmt.Errorf("Invalid external ID %d for request ID %d", externalID, requestID)
	}

	// Verify validator hasn't reported on the request.
	reported := false
	for _, report := range request.Reports {
		if report.Validator.Equals(validator) {
			reported = true
			break
		}
	}

	if reported {
		return nil, 0, fmt.Errorf(
			"Validator %s already submitted data report for this request", validator,
		)
	}

	// Verify request has not been expired
	params, _, err := queryParams(route, cliCtx)
	if err != nil {
		return nil, 0, err
	}

	if request.Request.RequestHeight+int64(params.ExpirationBlockCount) < height {
		return nil, 0, fmt.Errorf("Request #%d is already expired", requestID)
	}
	bz, err := types.QueryOK(VerificationResult{
		ChainID:      chainID,
		Validator:    validator,
		RequestID:    requestID,
		ExternalID:   externalID,
		DataSourceID: dataSourceID,
	})
	return bz, height, err
}
