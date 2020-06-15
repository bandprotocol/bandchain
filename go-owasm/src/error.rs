#[repr(i32)]
#[derive(Debug, PartialEq)]
pub enum Error {
  NoError = 0,
  CompliationError = 1,
  RunError = 2,
  ParseError = 3,
  WriteBinaryError = 4,
  ResolveNamesError = 5,
  ValidateError = 6,
  UnknownError = 7,
  SpanExceededCapacityError = 8
}
