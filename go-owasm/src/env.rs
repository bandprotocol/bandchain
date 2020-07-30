use crate::error::Error;
use crate::span::Span;

#[repr(C)]
pub struct env_t {
    _private: [u8; 0],
}

#[repr(C)]
// A struct representing the set of functions Rust can call back to Golang.
pub struct EnvDispatcher {
    pub get_calldata: extern "C" fn(*mut env_t, calldata: &mut Span) -> Error,
    pub set_return_data: extern "C" fn(*mut env_t, data: Span) -> Error,
    pub get_ask_count: extern "C" fn(*mut env_t) -> i64,
    pub get_min_count: extern "C" fn(*mut env_t) -> i64,
    pub get_ans_count: extern "C" fn(*mut env_t, &mut i64) -> Error,
    pub ask_external_data: extern "C" fn(*mut env_t, eid: i64, did: i64, data: Span) -> Error,
    pub get_external_data_status: extern "C" fn(*mut env_t, eid: i64, vid: i64, status: &mut i64) -> Error,
    pub get_external_data: extern "C" fn(*mut env_t, eid: i64, vid: i64, data: &mut Span) -> Error,
}

#[repr(C)]
// An execution environment passed from Golang world to Rust.
pub struct Env {
    pub env: *mut env_t,
    pub dis: EnvDispatcher,
}

#[repr(C)]
pub struct RunOutput {
    pub gas_used: u32,
}
