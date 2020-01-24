use owasm::core::{decode_outputs, encode_cmds};
use std::mem;

mod logic;

fn __return(output: &[u8]) -> u64 {
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
    __return(&bincode::config().big_endian().serialize(&logic::execute(outputs)).unwrap())
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

#[cfg(test)]
mod tests {
    use super::*;

    use owasm::ext::crypto::coins;

    #[test]
    fn test_encode_decode_parameter() {
        let params = logic::__Params { symbol: coins::Coins::ETH };

        let encoded_params = __encode_params(params).unwrap();
        let new_params: logic::Parameter =
            bincode::config().big_endian().deserialize(&encoded_params).ok().unwrap();

        println!("{:x?}", encoded_params);
        assert_eq!(new_params.symbol, coins::Coins::ETH);
    }

    // #[test]
    // fn test_prepare() {
    //     let params =
    //         logic::__Params { symbol_cg: String::from("ethereum"), symbol_cc: String::from("ETH") };
    //     let ptr = __return(&__encode_params(params).unwrap());

    //     let cmds = decode_cmds(__read_data(__prepare(ptr))).unwrap();
    //     assert_eq!(
    //         serde_json::to_string(&cmds[0]).unwrap(),
    //         r#"{"cmd":"curl","args":["https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"]}"#
    //     );
    //     assert_eq!(
    //         serde_json::to_string(&cmds[1]).unwrap(),
    //         r#"{"cmd":"curl","args":["https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD"]}"#
    //     );
    // }

    // #[test]
    // fn test_name() {
    //     let ptr_output = __name();
    //     println!("ptr is {}", ptr_output);
    //     let sl =
    //         unsafe { std::slice::from_raw_parts((ptr_output & ((1 << 32) - 1)) as *const u8, 12) };
    //     println!("length {}", sl.len());
    //     println!("x {}", sl[0]);
    //     // let result = std::str::from_utf8(__read_data(ptr_output)).unwrap();
    //     let result = std::str::from_utf8(sl).unwrap();
    //     assert_eq!(result, "Crypto price");
    // }

    // #[test]
    // fn test_parse_params() {
    //     let params =
    //         logic::__Params { symbol_cg: String::from("ethereum"), symbol_cc: String::from("ETH") };
    //     let encoded_params = __encode_params(params).unwrap();
    //     let ptr = __return(&encoded_params);

    //     let ptr_output = __parse_params(ptr);

    //     let result = std::str::from_utf8(__read_data(ptr_output)).unwrap();
    //     assert_eq!(result, r#"{"symbol_cg":"ethereum","symbol_cc":"ETH"}"#);
    // }

    // #[test]
    // fn test_raw_data_info() {
    //     let ptr_output = __raw_data_info();
    //     let result = std::str::from_utf8(__read_data(ptr_output)).unwrap();
    //     assert_eq!(result, r#"[["coin_gecko","f32"],["crypto_compare","f32"]]"#);
    // }

    // #[test]
    // fn test_params_info() {
    //     let ptr_output = __params_info();
    //     let result = std::str::from_utf8(__read_data(ptr_output)).unwrap();
    //     assert_eq!(result, r#"[["symbol_cg","String"],["symbol_cc","String"]]"#);
    // }
}
