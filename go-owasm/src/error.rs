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
  UnknownError = 255,
}
