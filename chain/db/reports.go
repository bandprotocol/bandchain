package db

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

func (b *BandDB) handleMsgReportData(
	txHash []byte,
	msg oracle.MsgReportData,
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

	for _, data := range msg.DataSet {
		// rawDataRequest, errSdk := b.OracleKeeper.GetRawRequest(
		// 	b.ctx, msg.RequestID, data.ExternalID,
		// )
		err := b.tx.Create(&ReportDetail{
			RequestID:  int64(msg.RequestID),
			Validator:  msg.Validator.String(),
			ExternalID: int64(data.ExternalID),
			// DataSourceID: 0, // TODO: FIX ME: Remove this col. Frontend can dig this itself.
			Data:     data.Data,
			Exitcode: data.ExitCode,
		}).Error

		if err != nil {
			return err
		}
	}

	return nil
}
