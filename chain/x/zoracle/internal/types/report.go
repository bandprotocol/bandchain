package types

// // Report represent detail of validator report
// type Report struct {
// 	Data       []RawDataReport `json:"data"`
// 	ReportedAt int64           `json:"reportedAt"`
// }

// // NewReport is a contructor of Report
// func NewReport(data []RawDataReport, reportedAt int64) Report {
// 	return Report{
// 		Data:       data,
// 		ReportedAt: reportedAt,
// 	}
// }

// RawDataReport encapsulates a raw data report for an external data source from a block validator.
type RawDataReport struct {
	ExternalDataID int64  `json:"externalDataID"`
	Data           []byte `json:"data"`
}

// NewRawDataReport creates a new RawDataReport instance.
func NewRawDataReport(externalDataID int64, data []byte) RawDataReport {
	return RawDataReport{
		ExternalDataID: externalDataID,
		Data:           data,
	}
}

// // ReportWithValidator is a report that contain operator address in struct
// type ReportWithValidator struct {
// 	Data       []RawDataReport `json:"data"`
// 	ReportedAt int64           `json:"reportedAt"`
// 	Validator  sdk.ValAddress  `json:"validator"`
// }

// // NewReportWithValidator is a contructor of ReportWithValidator
// func NewReportWithValidator(
// 	data []RawDataReport,
// 	reportedAt int64,
// 	valAddress sdk.ValAddress,
// ) ReportWithValidator {
// 	return ReportWithValidator{
// 		Data:       data,
// 		ReportedAt: reportedAt,
// 		Validator:  valAddress,
// 	}
// }
