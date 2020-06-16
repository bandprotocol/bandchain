#[repr(i32)]
#[derive(Debug, PartialEq, Clone)]
pub enum Error {
  NoError = 0,
  CompliationError = 1,
  RunError = 2,
  ParseError = 3,
  WriteBinaryError = 4,
  ResolveNamesError = 5,
  ValidateError = 6,
  SpanExceededCapacityError = 7,
  DeserializationError = 8,
  GasCounterInjectionError = 9,
  SerializationError = 10,
  GasLimitExceedError = 11,
  NoMemoryWasmError = 12,
  MinimumMemoryExceedError = 13,
  SetMaximumMemoryError = 14,
  StackHeightInstrumentationError = 15,
  CheckWasmImportsError = 16,
  CheckWasmExportsError = 17,
  InvalidSignatureFunctionError = 18,
  UnknownError = 255,
}

// TODO: Declare more error to cover exec env error
#[repr(i32)]
pub enum GoResult {
  Ok = 0,
  SetReturnDataWrongPeriod = 1,
  AnsCountWrongPeriod = 2,
  AskExternalDataWrongPeriod = 3,
  GetExternalDataStatusWrongPeriod = 4,
  GetExternalDataWrongPeriod = 5,
  GetExternalDataFromUnreportedValidator = 6,
  SpanExceededCapacity = 7,
  /// An error happened during normal operation of a Go callback
  Other = 8,
}

impl From<GoResult> for Error {
  fn from(r: GoResult) -> Self {
    match r {
      GoResult::Ok => Error::NoError,
      GoResult::SetReturnDataWrongPeriod => Error::InvalidFunctionCall,
      GoResult::AnsCountWrongPeriod => Error::InvalidFunctionCall,
      GoResult::AskExternalDataWrongPeriod => Error::InvalidFunctionCall,
      GoResult::GetExternalDataStatusWrongPeriod => Error::InvalidFunctionCall,
      GoResult::GetExternalDataWrongPeriod => Error::InvalidFunctionCall,
      GoResult::GetExternalDataFromUnreportedValidator => Error::InvalidFunctionCall,
      GoResult::SpanExceededCapacity => Error::SpanExceededCapacityError,
      GoResult::Other => Error::UnknownError,
    }
  }
}
