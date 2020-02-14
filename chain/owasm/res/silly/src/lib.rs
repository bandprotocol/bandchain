extern "C" {
    pub fn saveReturnData(dataOffset: *const u8, dataLength: i32);

    pub fn requestExternalData(
        dataSourceID: i64,
        externalDataID: i64,
        dataOffset: *const u8,
        dataLength: i32,
    );

    pub fn getExternalDataSize(externalDataID: i64, validatorIndex: i64) -> i32;

    pub fn readExternalData(
        externalDataID: i64,
        validatorIndex: i64,
        resultOffset: *mut u8,
        seekOffset: i32,
        resultSize: i32,
    );
}

#[no_mangle]
pub extern "C" fn prepare() {
    let message = "band-protocol".as_bytes();
    unsafe {
        requestExternalData(1, 1, message.as_ptr(), message.len() as i32);
    }
}

#[no_mangle]
pub extern "C" fn execute() {
    unsafe {
        let data_size = getExternalDataSize(1, 0);
        let mut data = vec![0u8; data_size as usize];
        readExternalData(1, 0, data.as_mut_ptr(), 0, data_size);
        saveReturnData(data.as_ptr(), data_size);
    }
}
