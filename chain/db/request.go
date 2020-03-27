package db

import (
	"strconv"

	"github.com/bandprotocol/bandchain/chain/x/zoracle"
)

func createRequest(
	id int64,
	oracleScriptID int64,
	calldata []byte,
	sufficientValidatorCount int64,
	expirationHeight int64,
	resolveStatus string,
	requester string,
	txHash []byte,
	result []byte,
) Request {
	return Request{
		ID:                       id,
		OracleScriptID:           oracleScriptID,
		Calldata:                 calldata,
		SufficientValidatorCount: sufficientValidatorCount,
		ExpirationHeight:         expirationHeight,
		ResolveStatus:            resolveStatus,
		Requester:                requester,
		TxHash:                   txHash,
		Result:                   result,
	}
}

func (b *BandDB) AddRequest(
	id int64,
	oracleScriptID int64,
	calldata []byte,
	sufficientValidatorCount int64,
	expirationHeight int64,
	resolveStatus string,
	requester string,
	txHash []byte,
	result []byte,
) error {
	request := createRequest(
		id,
		oracleScriptID,
		calldata,
		sufficientValidatorCount,
		expirationHeight,
		resolveStatus,
		requester,
		txHash,
		result,
	)
	err := b.tx.Create(&request).Error
	return err
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

func (b *BandDB) AddRawDataRequests(
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
	msg zoracle.MsgRequestData,
	events map[string]string,
) error {

	id, err := strconv.ParseInt(events[zoracle.EventTypeRequest+"."+zoracle.AttributeKeyID], 10, 64)
	if err != nil {
		return err
	}

	request := createRequest(
		id,
		int64(msg.OracleScriptID),
		msg.Calldata,
		msg.SufficientValidatorCount,
		msg.Expiration,
		"Pending",
		msg.Sender.String(),
		txHash,
		nil,
	)

	err = b.tx.Save(&request).Error
	if err != nil {
		return err
	}

	req, err := b.ZoracleKeeper.GetRequest(b.ctx, zoracle.RequestID(id))
	if err != nil {
		return err
	}

	for _, validatorAddress := range req.RequestedValidators {
		requestedValidator := createRequestedValidator(id, validatorAddress.String())
		err = b.tx.Save(&requestedValidator).Error
		if err != nil {
			return err
		}
	}

	for _, raw := range b.ZoracleKeeper.GetRawDataRequestWithExternalIDs(b.ctx, zoracle.RequestID(id)) {
		rawDataRequests := createRawDataRequests(id, int64(raw.ExternalID), int64(raw.RawDataRequest.DataSourceID), raw.RawDataRequest.Calldata)
		err = b.tx.Save(&rawDataRequests).Error
		if err != nil {
			return err
		}

		b.tx.FirstOrCreate(&RelatedDataSources{
			DataSourceID:   int64(raw.RawDataRequest.DataSourceID),
			OracleScriptID: int64(msg.OracleScriptID),
		})

	}

	return nil
}
