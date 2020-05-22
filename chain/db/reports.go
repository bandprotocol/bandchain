package db

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle"
)

func (b *BandDB) handleMsgReportData(
	txHash []byte,
	msg oracle.MsgReportData,
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
		var rawDataRequest RawDataRequests
		err := b.tx.Where(&RawDataRequests{
			RequestID:  int64(msg.RequestID),
			ExternalID: int64(data.ExternalID),
		}).First(&rawDataRequest).Error
		if err != nil {
			return err
		}

		err = b.tx.Create(&ReportDetail{
			RequestID:    int64(msg.RequestID),
			Validator:    msg.Validator.String(),
			ExternalID:   int64(data.ExternalID),
			DataSourceID: rawDataRequest.DataSourceID,
			Data:         data.Data,
			Exitcode:     data.ExitCode,
		}).Error

		if err != nil {
			return err
		}
	}

	return nil
}
