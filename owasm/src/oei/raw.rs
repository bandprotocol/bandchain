extern "C" {
    pub fn getCurrentRequestID() -> i64;
    pub fn getRequestedValidatorCount() -> i64;
    pub fn getSufficientValidatorCount() -> i64;
    pub fn getReceivedValidatorCount() -> i64;
    pub fn getPrepareBlockTime() -> i64;
    pub fn getAggregateBlockTime() -> i64;
    pub fn readValidatorAddress(validatorIndex: i64, resultOffset: *mut u8) -> i64;
    pub fn getCallDataSize() -> i64;
    pub fn readCallData(resultOffset: *mut u8, seekOffset: i64, resultSize: i64) -> i64;
    pub fn saveReturnData(dataOffset: *const u8, dataLength: i64) -> i64;
    pub fn requestExternalData(
        dataSourceID: i64,
        externalDataID: i64,
        dataOffset: *const u8,
        dataLength: i64,
    ) -> i64;
    pub fn getExternalDataSize(externalDataID: i64, validatorIndex: i64) -> i64;
    pub fn readExternalData(
        externalDataID: i64,
        validatorIndex: i64,
        resultOffset: *mut u8,
        seekOffset: i64,
        resultSize: i64,
    ) -> i64;
}
