#[repr(i32)]
#[derive(Debug, PartialEq, Clone)]
// An enum representing all kinds of errors we have in the system, with 0 for no error.
pub enum Error {
  NoError = 0,
  SpanTooSmallError = 1, // Span to write is too small.
  // Rust-generated errors during compilation.
  ValidationError = 2,           // Wasm code does not pass basic validation.
  DeserializationError = 3,      // Fail to deserialize Wasm into Partity-wasm module.
  SerializationError = 4,        // Fail to serialize Parity-wasm module into Wasm.
  InvalidImportsError = 5,       // Wasm code contains invalid import symbols.
  InvalidExportsError = 6,       // Wasm code contains invalid export symbols.
  BadMemorySectionError = 7,     // Wasm code contains bad memory sections.
  GasCounterInjectionError = 8,  // Fail to inject gas counter into Wasm code.
  StackHeightInjectionError = 9, // Fail to inject stack height limit into Wasm code.
  // Rust-generated errors during runtime.
  InstantiationError = 10,     // Error while instantiating Wasm with resolvers.
  RuntimeError = 11,           // Runtime error while executing the Wasm script.
  OutOfGasError = 12,          // Out-of-gas while executing the Wasm script.
  BadEntrySignatureError = 13, // Bad execution entry point sigature.
  // Go-generated errors while interacting with OEI.
  WrongPeriodActionError = 128,       // OEI action to invoke is not available.
  TooManyExternalDataError = 129,     // Too many external data requests.
  BadValidatorIndexError = 130,       // Bad validator index parameter.
  BadExternalIDError = 131,           // Bad external ID parameter.
  UnavailableExternalDataError = 132, // External data is not available.
  // Unexpected error
  UnknownError = 255,
}
