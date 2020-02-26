use owasm::oei;
use std::convert::TryInto;

#[no_mangle]
pub fn prepare() {}

#[no_mangle]
pub fn execute() {
    let sz = u64::from_le_bytes(
        TryInto::<[u8; 8]>::try_into((&oei::get_calldata()) as &[u8]).unwrap_or([0; 8]),
    );
    Vec::<i64>::with_capacity(sz as usize);
}
