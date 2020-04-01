package db

import (
	"github.com/bandprotocol/bandchain/chain/x/zoracle"
)

func (b *BandDB) handleMsgReportData(
	txHash []byte,
	msg zoracle.MsgReportData,
	events map[string]string,
) error {

	err := b.tx.Create(&Report{
		RequestID: int64(msg.RequestID),
		Validator: msg.Validator.String(),
		TxHash:    txHash,
		Reporter:  msg.Reporter.String(),
	}).Error

	if err != nil {
		return err
	}

	// for _, data := range msg.DataSet {
	// 	rawDataRequest, errSdk := b.ZoracleKeeper.GetRawDataRequest(
	// 		b.ctx, msg.RequestID, data.ExternalDataID,
	// 	)
	// 	if errSdk != nil {
	// 		return errSdk
	// 	}
	// 	err := b.tx.Create(&ReportDetail{
	// 		RequestID:    int64(msg.RequestID),
	// 		Validator:    msg.Validator.String(),
	// 		ExternalID:   int64(data.ExternalDataID),
	// 		DataSourceID: int64(rawDataRequest.DataSourceID),
	// 		Data:         data.Data,
	// 		Exitcode:     data.ExitCode,
	// 	}).Error

	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
