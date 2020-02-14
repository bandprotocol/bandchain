use std::mem::transmute;

extern "C" {
    pub fn getCurrentRequestID() -> i64;
    pub fn getRequestedValidatorCount() -> i64;

    pub fn saveReturnData(dataOffset: *const u8, dataLength: i32);
}

#[no_mangle]
pub extern "C" fn execute() {
    unsafe {
        let req_id = getCurrentRequestID();
        let validators_count = getRequestedValidatorCount();

        let arr: [u8; 8] = transmute((req_id * 10 + validators_count).to_be());
        let mut data = Box::new(arr).to_vec();
        data.reverse();
        saveReturnData(data.as_ptr(), 8);
    }
}
