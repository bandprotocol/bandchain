use owasm::core::{decode_cmds, decode_outputs, encode_cmds};
use std::mem;

mod logic;

fn __return(output: &[u8]) -> *const u8 {
    let sz = output.len();
    let loc = __allocate(4 + sz as usize);
    unsafe {
        std::ptr::copy_nonoverlapping(&(sz as u32), loc as *mut u32, 4);
        std::ptr::copy_nonoverlapping(output.as_ptr(), loc.offset(4), sz)
    };
    loc
}

/// Encodes parameter struct to `Vec<u8>` for testing only.
fn __encode_params(params: logic::__Params) -> Option<Vec<u8>> {
    bincode::config().big_endian().serialize(&params).ok()
}

fn __decode_params(input: *const u8) -> Option<logic::__Params> {
    // Get raw params from memory
    let data =
        unsafe { std::slice::from_raw_parts(input.offset(4), *(input as *const u32) as usize) };
    bincode::config().big_endian().deserialize(data).ok()
}

fn __decode_data(params: &logic::__Params, input: *const u8) -> Option<logic::__Data> {
    unsafe {
        let each_size = *(input as *const u32);
        logic::__Data::__output(
            &params,
            decode_outputs(std::slice::from_raw_parts(input.offset(4), each_size as usize))?,
        )
    }
}

#[no_mangle]
pub fn __allocate(size: usize) -> *mut u8 {
    let mut buffer = Vec::with_capacity(size);
    let pointer = buffer.as_mut_ptr();
    mem::forget(buffer);
    pointer
}

#[no_mangle]
pub fn __prepare(params: *const u8) -> *const u8 {
    __return(&encode_cmds(logic::__Data::__input(&__decode_params(params).unwrap())).unwrap())
}

#[no_mangle]
pub fn __execute(params: *const u8, input: *const *const u8) -> *const u8 {
    let p = __decode_params(params).unwrap();
    let outputs: Vec<_> = unsafe {
        std::slice::from_raw_parts(input.offset(1), *(input as *const u32) as usize)
            .to_vec()
            .into_iter()
            .filter_map(|each| __decode_data(&p, each))
            .collect()
    };
    __return(&bincode::config().big_endian().serialize(&logic::execute(outputs)).unwrap())
}

#[no_mangle]
pub fn __name() -> *const u8 {
    __return(&logic::name().into_bytes())
}

#[no_mangle]
pub fn __params_info() -> *const u8 {
    __return(&serde_json::to_string(&logic::__Params::__fields()).ok().unwrap().into_bytes())
}

#[no_mangle]
pub fn __parse_params(params: *const u8) -> *const u8 {
    __return(&serde_json::to_string(&__decode_params(params).unwrap()).ok().unwrap().into_bytes())
}

#[no_mangle]
pub fn __raw_data_info() -> *const u8 {
    __return(&serde_json::to_string(&logic::__Data::__fields()).ok().unwrap().into_bytes())
}

#[no_mangle]
pub fn __parse_raw_data(params: *const u8, input: *const u8) -> *const u8 {
    __return(
        &serde_json::to_string(&__decode_data(&__decode_params(params).unwrap(), input).unwrap())
            .ok()
            .unwrap()
            .into_bytes(),
    )
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_encode_decode_parameter() {
        let params =
            logic::__Params { symbol_cg: String::from("ethereum"), symbol_cc: String::from("ETH") };

        let encoded_params = __encode_params(params).unwrap();
        let ptr = __return(&encoded_params);

        let new_params = __decode_params(ptr).unwrap();

        assert_eq!(new_params.symbol_cg, String::from("ethereum"));
    }

    #[test]
    fn test_prepare() {
        let params =
            logic::__Params { symbol_cg: String::from("ethereum"), symbol_cc: String::from("ETH") };
        let ptr = __return(&__encode_params(params).unwrap());

        let prepare_ptr = __prepare(ptr);
        let data = unsafe {
            std::slice::from_raw_parts(prepare_ptr.offset(4), *(prepare_ptr as *const u32) as usize)
        };
        let cmds = decode_cmds(&data).unwrap();
        assert_eq!(
            serde_json::to_string(&cmds[0]).unwrap(),
            r#"{"cmd":"curl","args":["https://api.coingecko.com/api/v3/simple/price?ids=ethereum&vs_currencies=usd"]}"#
        );
        assert_eq!(
            serde_json::to_string(&cmds[1]).unwrap(),
            r#"{"cmd":"curl","args":["https://min-api.cryptocompare.com/data/price?fsym=ETH&tsyms=USD"]}"#
        );
    }

    #[test]
    fn test_name() {
        let ptr_output = __name();
        let data = unsafe {
            std::slice::from_raw_parts(ptr_output.offset(4), *(ptr_output as *const u32) as usize)
        };

        let result = std::str::from_utf8(data).unwrap();
        assert_eq!(result, "Crypto price");
    }

    #[test]
    fn test_parse_params() {
        let params =
            logic::__Params { symbol_cg: String::from("ethereum"), symbol_cc: String::from("ETH") };
        let encoded_params = __encode_params(params).unwrap();
        let ptr = __return(&encoded_params);

        let ptr_output = __parse_params(ptr);
        let data = unsafe {
            std::slice::from_raw_parts(ptr_output.offset(4), *(ptr_output as *const u32) as usize)
        };

        let result = std::str::from_utf8(data).unwrap();
        assert_eq!(result, r#"{"symbol_cg":"ethereum","symbol_cc":"ETH"}"#);
    }

    #[test]
    fn test_raw_data_info() {
        let ptr_output = __raw_data_info();
        let data = unsafe {
            std::slice::from_raw_parts(ptr_output.offset(4), *(ptr_output as *const u32) as usize)
        };

        let result = std::str::from_utf8(data).unwrap();
        assert_eq!(result, r#"[["coin_gecko","f32"],["crypto_compare","f32"]]"#);
    }

    #[test]
    fn test_params_info() {
        let ptr_output = __params_info();
        let data = unsafe {
            std::slice::from_raw_parts(ptr_output.offset(4), *(ptr_output as *const u32) as usize)
        };

        let result = std::str::from_utf8(data).unwrap();
        assert_eq!(result, r#"[["symbol_cg","String"],["symbol_cc","String"]]"#);
    }
}
