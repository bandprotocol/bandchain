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
	ErrInstantiationError         = errors.New("compile fail")
	ErrRuntimeError               = errors.New("run fail")
	ErrValidationError            = errors.New("validate fail")
	ErrGasLimitExceeded           = errors.New("gas limit exceeded")
	ErrNoMemoryWasm               = errors.New("no memory wasm")
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
func parseErrorFromC(code C.Error) error {
	switch code {
	case C.Error_NoError:
		return nil
	case C.Error_InstantiationError:
		return ErrCompliationFail
	case C.Error_RuntimeError:
		return ErrRunFail
	case C.Error_ValidationError:
		return ErrValidateFail
	case C.Error_SpanTooSmallError:
		return ErrSpanExceededCapacity
	case C.Error_DeserializationError:
		return ErrDeserializeFail
	case C.Error_GasCounterInjectionError:
		return ErrGasCounterInjectFail
	case C.Error_SerializationError:
		return ErrSerializetFail
	case C.Error_OutOfGasError:
		return ErrGasLimitExceeded
	case C.Error_BadMemorySectionError:
		return ErrNoMemoryWasm
	case C.Error_StackHeightInjectionError:
		return ErrStackHeightInstrumentation
	case C.Error_InvalidImportsError:
		return ErrCheckWasmImports
	case C.Error_InvalidExportsError:
		return ErrCheckWasmExports
	case C.Error_BadEntrySignatureError:
		return ErrInvalidSignatureFunction
	case C.Error_WrongPeriodActionError:
		return ErrSetReturnDataWrongPeriod
	case C.Error_TooManyExternalDataError:
		return ErrAskExternalDataExceed
	case C.Error_BadValidatorIndexError:
		return ErrValidatorOutOfRange
	case C.Error_BadExternalIDError:
		return ErrInvalidExternalID
	case C.Error_UnavailbleExternalDataError:
		return ErrGetUnreportedData
	default:
		return ErrUnknownError
	}
}
