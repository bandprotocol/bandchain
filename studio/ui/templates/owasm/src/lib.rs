use owasm::core::{decode_outputs, encode_cmds};
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

#[no_mangle]
pub fn __allocate(size: usize) -> *mut u8 {
    let mut buffer = Vec::with_capacity(size);
    let pointer = buffer.as_mut_ptr();
    mem::forget(buffer);
    pointer
}

#[no_mangle]
pub fn __prepare() -> *const u8 {
    __return(&encode_cmds(logic::__Data::__input()).unwrap())
}

#[no_mangle]
pub fn __execute(input: *const *const u8) -> *const u8 {
    let outputs: Vec<_> = unsafe {
        std::slice::from_raw_parts(input.offset(1), *(input as *const u32) as usize)
            .to_vec()
            .into_iter()
            .filter_map(|each| {
                let each_size = *(each as *const u32);
                logic::__Data::__output(decode_outputs(std::slice::from_raw_parts(
                    each.offset(4),
                    each_size as usize,
                ))?)
            })
            .collect()
    };
    __return(&bincode::config().big_endian().serialize(&logic::execute(outputs)).unwrap())
}
