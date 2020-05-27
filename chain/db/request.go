package db

import (
	"errors"
	"strconv"

	"github.com/bandprotocol/bandchain/chain/x/oracle"
	otypes "github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

const (
	Open    = "Pending"
	Success = "Success"
	Failure = "Failure"
	Expired = "Expired"
	Unknown = "Unknown"
)

func parseResolveStatus(resolveStatus otypes.ResolveStatus) string {
	switch resolveStatus {
	case 0:
		return Open
	case 1:
		return Success
	case 2:
		return Failure
	case 3:
		return Expired
	default:
		return Unknown
	}
}

func createRequest(
	id int64,
	oracleScriptID int64,
	calldata []byte,
	minCount int64,
	expirationHeight int64,
	resolveStatus string,
	requester string,
	clientID string,
	txHash []byte,
	result []byte,
) Request {
	return Request{
		ID:               id,
		OracleScriptID:   oracleScriptID,
		Calldata:         calldata,
		MinCount:         minCount,
		ExpirationHeight: expirationHeight,
		ResolveStatus:    resolveStatus,
		Requester:        requester,
		ClientID:         clientID,
		TxHash:           txHash,
		Result:           result,
	}
}

func (b *BandDB) AddNewRequest(
	id int64,
	oracleScriptID int64,
	calldata []byte,
	minCount int64,
	expirationHeight int64,
	resolveStatus string,
	requester string,
	clientID string,
	txHash []byte,
	result []byte,
	rawExternalIDs []string,
	rawDataSourceIDs []string,
	rawCalldatas []string,
) error {
	request := createRequest(
		id,
		oracleScriptID,
		calldata,
		minCount,
		expirationHeight,
		resolveStatus,
		requester,
		clientID,
		txHash,
		result,
	)
	err := b.tx.Create(&request).Error
	if err != nil {
		return err
	}

	req, err := b.OracleKeeper.GetRequest(b.ctx, otypes.RequestID(id))
	if err != nil {
		return err
	}

	for _, validatorAddress := range req.RequestedValidators {
		err := b.AddRequestedValidator(id, validatorAddress.String())
		if err != nil {
			return err
		}
	}

	for i := range rawExternalIDs {
		externalID, err := strconv.ParseInt(rawExternalIDs[i], 10, 64)
		if err != nil {
			return err
		}
		dataSourceID, err := strconv.ParseInt(rawDataSourceIDs[i], 10, 64)
		if err != nil {
			return err
		}
		err = b.AddRawDataRequest(
			id,
			externalID,
			dataSourceID,
			[]byte(rawCalldatas[i]),
		)
		if err != nil {
			return err
		}
		err = b.tx.FirstOrCreate(&RelatedDataSources{
			DataSourceID:   dataSourceID,
			OracleScriptID: oracleScriptID,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func createRequestedValidator(
	requestID int64,
	validatorAddress string,
) RequestedValidator {
	return RequestedValidator{
		RequestID:        requestID,
		ValidatorAddress: validatorAddress,
	}
}

func (b *BandDB) AddRequestedValidator(
	requestID int64,
	validatorAddress string,
) error {
	requestValidator := createRequestedValidator(
		requestID,
		validatorAddress,
	)
	err := b.tx.Create(&requestValidator).Error
	return err
}

func createRawDataRequests(
	requestID int64,
	externalID int64,
	dataSourceID int64,
	calldata []byte,
) RawDataRequests {
	return RawDataRequests{
		RequestID:    requestID,
		ExternalID:   externalID,
		DataSourceID: dataSourceID,
		Calldata:     calldata,
	}
}

func (b *BandDB) AddRawDataRequest(
	requestID int64,
	externalID int64,
	dataSourceID int64,
	calldata []byte,
) error {
	rawDataRequests := createRawDataRequests(
		requestID,
		externalID,
		dataSourceID,
		calldata,
	)
	err := b.tx.Create(&rawDataRequests).Error
	return err
}

func (b *BandDB) handleMsgRequestData(
	txHash []byte,
	msg oracle.MsgRequestData,
	events map[string][]string,
) error {
	ids := events[otypes.EventTypeRequest+"."+otypes.AttributeKeyID]
	if len(ids) != 1 {
		return errors.New("handleMsgCreateDataSource: cannot find request id")
	}
	id, err := strconv.ParseInt(ids[0], 10, 64)
	if err != nil {
		return err
	}
	request, err := b.OracleKeeper.GetRequest(b.ctx, otypes.RequestID(id))
	if err != nil {
		return err
	}
	return b.AddNewRequest(
		id,
		int64(msg.OracleScriptID),
		msg.Calldata,
		msg.MinCount,
		request.RequestHeight+20, // TODO: REMOVE THIS. HACK!!!
		"Pending",
		msg.Sender.String(),
		msg.ClientID,
		txHash,
		nil,
		events[otypes.EventTypeRawRequest+"."+otypes.AttributeKeyExternalID],
		events[otypes.EventTypeRawRequest+"."+otypes.AttributeKeyDataSourceID],
		events[otypes.EventTypeRawRequest+"."+otypes.AttributeKeyCalldata],
	)
}
