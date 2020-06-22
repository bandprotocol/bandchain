package api

// #include "bindings.h"
import "C"
import (
	"errors"
)

var (
	ErrSpanTooSmall            = errors.New("span to write is too small")
	ErrValidation              = errors.New("wasm code does not pass basic validation")
	ErrDeserialization         = errors.New("fail to deserialize wasm into Partity-wasm module")
	ErrSerialization           = errors.New("fail to serialize Parity-wasm module into wasm")
	ErrInvalidImports          = errors.New("wasm code contains invalid import symbols")
	ErrInvalidExports          = errors.New("wasm code contains invalid export symbols")
	ErrBadMemorySection        = errors.New("wasm code contains bad memory sections")
	ErrGasCounterInjection     = errors.New("fail to inject gas counter into Wasm code")
	ErrStackHeightInjection    = errors.New("fail to inject stack height limit into Wasm code")
	ErrInstantiation           = errors.New("error while instantiating Wasm with resolvers")
	ErrRuntime                 = errors.New("runtime error while executing the Wasm script")
	ErrOutOfGas                = errors.New("out-of-gas while executing the wasm script")
	ErrBadEntrySignature       = errors.New("bad execution entry point sigature")
	ErrMemoryOutOfBound        = errors.New("out-of-bound memory access while executing the wasm script")
	ErrWrongPeriodAction       = errors.New("OEI action to invoke is not available")
	ErrTooManyExternalData     = errors.New("too many external data requests")
	ErrBadValidatorIndex       = errors.New("bad validator index parameter")
	ErrBadExternalID           = errors.New("bad external ID parameter")
	ErrUnavailableExternalData = errors.New("external data is not available")
	ErrUnknown                 = errors.New("unknown error")
)

// toCError converts the given Golang error into Rust/C error.
func toCError(err error) C.Error {
	switch err {
	case nil:
		return C.Error_NoError
	case ErrWrongPeriodAction:
		return C.Error_WrongPeriodActionError
	case ErrTooManyExternalData:
		return C.Error_TooManyExternalDataError
	case ErrBadValidatorIndex:
		return C.Error_BadValidatorIndexError
	case ErrBadExternalID:
		return C.Error_BadExternalIDError
	case ErrUnavailableExternalData:
		return C.Error_UnavailableExternalDataError
	default:
		return C.Error_UnknownError
	}
}

// toGoError converts the given Rust/C error to Golang error.
func toGoError(code C.Error) error {
	switch code {
	case C.Error_NoError:
		return nil
	case C.Error_SpanTooSmallError:
		return ErrSpanTooSmall
	// Rust-generated errors during compilation.
	case C.Error_ValidationError:
		return ErrValidation
	case C.Error_DeserializationError:
		return ErrDeserialization
	case C.Error_SerializationError:
		return ErrSerialization
	case C.Error_InvalidImportsError:
		return ErrInvalidImports
	case C.Error_InvalidExportsError:
		return ErrInvalidExports
	case C.Error_BadMemorySectionError:
		return ErrBadMemorySection
	case C.Error_GasCounterInjectionError:
		return ErrGasCounterInjection
	case C.Error_StackHeightInjectionError:
		return ErrStackHeightInjection
	// Rust-generated errors during runtime.
	case C.Error_InstantiationError:
		return ErrInstantiation
	case C.Error_RuntimeError:
		return ErrRuntime
	case C.Error_OutOfGasError:
		return ErrOutOfGas
	case C.Error_BadEntrySignatureError:
		return ErrBadEntrySignature
	// Go-generated errors while interacting with OEI.
	case C.Error_WrongPeriodActionError:
		return ErrWrongPeriodAction
	case C.Error_TooManyExternalDataError:
		return ErrTooManyExternalData
	case C.Error_BadValidatorIndexError:
		return ErrBadValidatorIndex
	case C.Error_BadExternalIDError:
		return ErrBadExternalID
	case C.Error_UnavailableExternalDataError:
		return ErrUnavailableExternalData
	default:
		return ErrUnknown
	}
}
