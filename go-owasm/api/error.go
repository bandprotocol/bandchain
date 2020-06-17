package api

// #include "bindings.h"
import "C"
import (
	"errors"
)

var (
	ErrCompliationFail            = errors.New("compile fail")
	ErrRunFail                    = errors.New("run fail")
	ErrParseFail                  = errors.New("parse fail")
	ErrWriteBinaryFail            = errors.New("write binary fail")
	ErrResolvesNamesFail          = errors.New("resolve names fail")
	ErrValidateFail               = errors.New("validate fail")
	ErrSpanExceededCapacity       = errors.New("span exceeded capacity")
	ErrDeserializeFail            = errors.New("deserialize fail")
	ErrGasCounterInjectFail       = errors.New("gas counter inject fail")
	ErrSerializetFail             = errors.New("serialize fail")
	ErrUnknownError               = errors.New("unknown error")
	ErrCompliationError           = errors.New("compile fail")
	ErrRunError                   = errors.New("run fail")
	ErrParseError                 = errors.New("parse fail")
	ErrWriteBinaryError           = errors.New("write binary fail")
	ErrValidateError              = errors.New("validate fail")
	ErrGasLimitExceeded           = errors.New("gas limit exceeded")
	ErrNoMemoryWasm               = errors.New("no memory wasm")
	ErrMinimumMemoryExceed        = errors.New("minimum memory exceed")
	ErrSetMaximumMemory           = errors.New("maximum must be unset")
	ErrStackHeightInstrumentation = errors.New("stack height instrumention fail")
	ErrCheckWasmImports           = errors.New("check wasm imports fail")
	ErrCheckWasmExports           = errors.New("check wasm exports fail")
	ErrInvalidSignatureFunction   = errors.New("invalid signature function")

	ErrSetReturnDataWrongPeriod         = errors.New("set return data on non-execution period")
	ErrAnsCountWrongPeriod              = errors.New("get ans count on non-execution period")
	ErrAskExternalDataWrongPeriod       = errors.New("ask external data on non-prepare period")
	ErrAskExternalDataExceed            = errors.New("ask external data exceed")
	ErrGetExternalDataStatusWrongPeriod = errors.New("get external data status on non-execution period")
	ErrGetExternalDataWrongPeriod       = errors.New("get external data on non-execution period")
	ErrValidatorOutOfRange              = errors.New("validator index out of range")
	ErrInvalidExternalID                = errors.New("get data from invalid external id")
	ErrGetUnreportedData                = errors.New("get data from unreported validator")
)

// parseError - returns parsed error from errors code on bindings.h
func parseError(code C.Error) error {
	switch code {
	case C.Error_NoError:
		return nil
	case C.Error_CompliationError:
		return ErrCompliationFail
	case C.Error_RunError:
		return ErrRunFail
	case C.Error_ParseError:
		return ErrParseFail
	case C.Error_WriteBinaryError:
		return ErrWriteBinaryFail
	case C.Error_ResolveNamesError:
		return ErrResolvesNamesFail
	case C.Error_ValidateError:
		return ErrValidateFail
	case C.Error_SpanExceededCapacityError:
		return ErrSpanExceededCapacity
	case C.Error_DeserializationError:
		return ErrDeserializeFail
	case C.Error_GasCounterInjectionError:
		return ErrGasCounterInjectFail
	case C.Error_SerializationError:
		return ErrSerializetFail
	case C.Error_InvalidSignatureFunctionError:
		return ErrInvalidSignatureFunction
	case C.Error_GasLimitExceedError:
		return ErrGasLimitExceeded
	case C.Error_NoMemoryWasmError:
		return ErrNoMemoryWasm
	case C.Error_MinimumMemoryExceedError:
		return ErrMinimumMemoryExceed
	case C.Error_SetMaximumMemoryError:
		return ErrSetMaximumMemory
	case C.Error_StackHeightInstrumentationError:
		return ErrStackHeightInstrumentation
	case 16:
		return ErrCheckWasmImports
	case 17:
		return ErrCheckWasmExports
	case 18:
		return ErrInvalidSignatureFunction
	case C.Error_SetReturnDataWrongPeriodError:
		return ErrSetReturnDataWrongPeriod
	case C.Error_AnsCountWrongPeriodError:
		return ErrAnsCountWrongPeriod
	case C.Error_AskExternalDataWrongPeriodError:
		return ErrAskExternalDataWrongPeriod
	case C.Error_AskExternalDataExceedError:
		return ErrAskExternalDataExceed
	case C.Error_GetExternalDataStatusWrongPeriodError:
		return ErrGetExternalDataStatusWrongPeriod
	case C.Error_GetExternalDataWrongPeriodError:
		return ErrGetExternalDataWrongPeriod
	case C.Error_ValidatorOutOfRangeError:
		return ErrValidatorOutOfRange
	case C.Error_InvalidExternalIDError:
		return ErrInvalidExternalID
	case C.Error_GetUnreportedDataError:
		return ErrGetUnreportedData
	default:
		return ErrUnknownError
	}
}
