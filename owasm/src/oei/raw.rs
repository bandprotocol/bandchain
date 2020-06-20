extern "C" {
    // pub fn getCurrentRequestID() -> i64;
    pub fn get_ask_count() -> i64;
    pub fn get_min_count() -> i64;
    pub fn get_ans_count() -> i64;
    // pub fn getPrepareBlockTime() -> i64;
    // pub fn getAggregateBlockTime() -> i64;
    // pub fn readValidatorAddress(validatorIndex: i64, resultOffset: *mut u8) -> i64;

    // TODO: Remove get_calldata_size and change signature of read_calldata to
    // read_calldata(resOffset) -> i64 output means lenght of calldata if calldata length exceed
    // limit or error occured between copy data to pointer it will return -1
    pub fn get_calldata_size() -> i64;
    pub fn read_calldata(resOffset: i64, resSize: i64);
    pub fn set_return_data(dataOffset: i64, dataLength: i64);
    pub fn ask_external_data(eid: i64, did: i64, dataOffset: i64, dataLength: i64);
    // TODO: Same as read_calldata
    pub fn get_external_data_size(eid: i64, vid: i64) -> i64;
    pub fn read_external_data(eid: i64, vid: i64, resOffset: i64, resSize: i64);
}
