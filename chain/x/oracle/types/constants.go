package types

// nolint
const (
	DoNotModify = "[do-not-modify]"

	MaxNameLength        = 128
	MaxDescriptionLength = 4096
	MaxClientIDLength    = 128
	MaxSchemaLength      = 512
	MaxURLLength         = 128

	MaxExecutableSize     = 8 * 1024   // 8kB
	MaxWasmCodeSize       = 512 * 1024 // 512kB
	MaxCalldataSize       = 1 * 1024   // 1kB
	MaxRawRequestDataSize = 1 * 1024   // 1kB
	MaxRawReportDataSize  = 1 * 1024   // 1kB
	MaxResultSize         = 1 * 1024   // 1kB

	WasmPrepareGas = 100000
	WasmExecuteGas = 100000
)

// nolint
var (
	DoNotModifyBytes = []byte(DoNotModify)
)
