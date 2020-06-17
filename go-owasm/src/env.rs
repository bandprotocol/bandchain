use crate::error::GoResult;
use crate::span::Span;

#[repr(C)]
pub struct env_t {
    _private: [u8; 0],
}

#[repr(C)]
pub struct EnvDispatcher {
    pub get_calldata: extern "C" fn(*mut env_t, calldata: &mut Span) -> GoResult,
    pub set_return_data: extern "C" fn(*mut env_t, data: Span) -> GoResult,
    pub get_ask_count: extern "C" fn(*mut env_t) -> i64,
    pub get_min_count: extern "C" fn(*mut env_t) -> i64,
    pub get_ans_count: extern "C" fn(*mut env_t, &mut i64) -> GoResult,
    pub ask_external_data: extern "C" fn(*mut env_t, eid: i64, did: i64, data: Span) -> GoResult,
    pub get_external_data_status:
        extern "C" fn(*mut env_t, eid: i64, vid: i64, status: &mut i64) -> GoResult,
    pub get_external_data:
        extern "C" fn(*mut env_t, eid: i64, vid: i64, data: &mut Span) -> GoResult,
}

#[repr(C)]
pub struct Env {
    pub env: *mut env_t,
    pub dis: EnvDispatcher,
}
