#[repr(i32)]
pub enum Error {
  NoError = 0,
  CompliationError = 1,
  RunError = 2,
  ParseError = 3,
  WriteBinaryError = 4,
  ResolveNamesError = 5,
  ValidateError = 6,
  UnknownError = 7,
  DeserializationError = 9,
  GasCounterInjectionError = 10,
  SerializationError = 11,
}
