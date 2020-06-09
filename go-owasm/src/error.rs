#[repr(i32)]
pub enum Error {
  Ok = 0,
  CompileFail = 1,
  RunFail = 2,
  ParseFail = 3,
  WriteBinaryFail = 4,
  ResolveNamesFail = 5,
  ValidateFail = 6,
  UnknownFail = 7
}
