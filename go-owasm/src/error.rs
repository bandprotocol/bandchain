#[repr(i32)]
pub enum Error {
  Ok = 0,
  CompileFail = 1,
  RunFail = 2,
  ParseFail = 3,
  Nul = 4,
  NonUtf8Result = 5
}
