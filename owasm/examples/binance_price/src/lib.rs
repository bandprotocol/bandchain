use owasm::core::{decode_outputs, encode_cmds};
use std::mem;

mod logic;

fn __return(output: &[u8]) -> u64 {
    /// TESTs
    let sz = output.len();
    let loc = __allocate(sz);
    unsafe { std::ptr::copy_nonoverlapping(output.as_ptr(), loc, sz) };
    ((sz as u64) << 32) | ((loc as u32) as u64)
}

fn __read_data<'a, T>(ptr: u64) -> &'a [T] {
    unsafe { std::slice::from_raw_parts((ptr & ((1 << 32) - 1)) as *const T, (ptr >> 32) as usize) }
}

/// Encodes parameter struct to `Vec<u8>` for testing only.
fn __encode_params(params: logic::__Params) -> Option<Vec<u8>> {
    bincode::config().big_endian().serialize(&params).ok()
}

fn __decode_params(input: u64) -> Option<logic::__Params> {
    bincode::config().big_endian().deserialize(__read_data(input)).ok()
}

fn __decode_data(params: &logic::__Params, input: u64) -> Option<logic::__Data> {
    logic::__Data::__output(&params, decode_outputs(__read_data(input))?)
}

fn __decode_result(input: u64) -> Option<logic::__Result> {
    bincode::config().big_endian().deserialize(__read_data(input)).ok()
}

#[no_mangle]
pub fn __allocate(size: usize) -> *mut u8 {
    let mut buffer = Vec::with_capacity(size);
    let pointer = buffer.as_mut_ptr();
    mem::forget(buffer);
    pointer
}

#[no_mangle]
pub fn __prepare(params: u64) -> u64 {
    __return(&encode_cmds(logic::__Data::__input(&__decode_params(params).unwrap())).unwrap())
}

#[no_mangle]
pub fn __execute(params: u64, input: u64) -> u64 {
    let p = __decode_params(params).unwrap();
    let outputs: Vec<_> = __read_data::<u64>(input)
        .to_vec()
        .into_iter()
        .filter_map(|each| __decode_data(&p, each))
        .collect();
    __return(&bincode::config().big_endian().serialize(&logic::execute(p, outputs)).unwrap())
}

#[no_mangle]
pub fn __params_info() -> u64 {
    __return(&serde_json::to_string(&logic::__Params::__fields()).ok().unwrap().into_bytes())
}

#[no_mangle]
pub fn __parse_params(params: u64) -> u64 {
    __return(&serde_json::to_string(&__decode_params(params).unwrap()).ok().unwrap().into_bytes())
}

#[no_mangle]
pub fn __serialize_params(json_ptr: u64) -> u64 {
    let params: logic::__Params = serde_json::from_str(
        String::from_utf8(__read_data(json_ptr).to_vec()).ok().unwrap().as_str(),
    )
    .ok()
    .unwrap();
    __return(&__encode_params(params).unwrap())
}

#[no_mangle]
pub fn __raw_data_info() -> u64 {
    __return(&serde_json::to_string(&logic::__Data::__fields()).ok().unwrap().into_bytes())
}

#[no_mangle]
pub fn __parse_raw_data(params: u64, input: u64) -> u64 {
    __return(
        &serde_json::to_string(&__decode_data(&__decode_params(params).unwrap(), input).unwrap())
            .ok()
            .unwrap()
            .into_bytes(),
    )
}

#[no_mangle]
pub fn __result_info() -> u64 {
    __return(&serde_json::to_string(&logic::__Result::__fields()).ok().unwrap().into_bytes())
}

#[no_mangle]
pub fn __parse_result(result: u64) -> u64 {
    __return(&serde_json::to_string(&__decode_result(result).unwrap()).ok().unwrap().into_bytes())
}

