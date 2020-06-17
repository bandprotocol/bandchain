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
  FunctionNotFoundError = 11,
  GasLimitExceedError = 12,
  NoMemoryWasmError = 13,
  MinimumMemoryExceedError = 14,
  SetMaximumMemoryError = 15,
  StackHeightInstrumentation = 16,
  UnknownError = 255,
}
