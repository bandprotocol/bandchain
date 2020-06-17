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

  SetReturnDataWrongPeriodError = 31,
  AnsCountWrongPeriodError = 32,
  AskExternalDataWrongPeriodError = 33,
  AskExternalDataExceedError = 34,
  GetExternalDataStatusWrongPeriodError = 35,
  GetExternalDataWrongPeriodError = 36,
  ValidatorOutOfRangeError = 37,
  InvalidExternalIDError = 38,
  GetUnreportedDataError = 39,

  UnknownError = 255,
}

#[repr(i32)]
pub enum GoResult {
  Ok = 0,
  SpanExceededCapacity = 1,
  SetReturnDataWrongPeriod = 2,
  AnsCountWrongPeriod = 3,
  AskExternalDataWrongPeriod = 4,
  AskExternalDataExceed = 5,
  GetExternalDataStatusWrongPeriod = 6,
  GetExternalDataWrongPeriod = 7,
  ValidatorOutOfRange = 8,
  InvalidExternalID = 9,
  GetUnreportedData = 10,
  Other = 11,
}

impl From<GoResult> for Error {
  fn from(r: GoResult) -> Self {
    match r {
      GoResult::Ok => Error::NoError,
      GoResult::SetReturnDataWrongPeriod => Error::SetReturnDataWrongPeriodError,
      GoResult::AnsCountWrongPeriod => Error::AnsCountWrongPeriodError,
      GoResult::AskExternalDataWrongPeriod => Error::AskExternalDataWrongPeriodError,
      GoResult::AskExternalDataExceed => Error::AskExternalDataExceedError,
      GoResult::GetExternalDataStatusWrongPeriod => Error::GetExternalDataStatusWrongPeriodError,
      GoResult::GetExternalDataWrongPeriod => Error::GetExternalDataWrongPeriodError,
      GoResult::ValidatorOutOfRange => Error::ValidatorOutOfRangeError,
      GoResult::InvalidExternalID => Error::InvalidExternalIDError,
      GoResult::GetUnreportedData => Error::GetUnreportedDataError,
      GoResult::SpanExceededCapacity => Error::SpanExceededCapacityError,
      GoResult::Other => Error::UnknownError,
    }
  }
}
