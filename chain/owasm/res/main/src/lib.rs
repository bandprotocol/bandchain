extern "C" {
    pub fn getCurrentRequestID() -> i64;
    pub fn getRequestedValidatorCount() -> i64;

    pub fn saveReturnData(dataOffset: *const u8, dataLength: i32);
}

#[no_mangle]
pub fn execute() {
    unsafe {
        let req_id = getCurrentRequestID();
        let validators_count = getRequestedValidatorCount();

        saveReturnData(
            vec![(req_id + validators_count) as u8, 0, 0, 0, 0, 0, 0, 0].as_ptr(),
            8,
        );
    }
}
