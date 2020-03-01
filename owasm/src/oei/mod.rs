mod raw;

pub fn get_request_id() -> i64 {
    unsafe { raw::getCurrentRequestID() }
}

pub fn get_requested_validator_count() -> i64 {
    unsafe { raw::getRequestedValidatorCount() }
}

pub fn get_sufficient_validator_count() -> i64 {
    unsafe { raw::getSufficientValidatorCount() }
}

pub fn get_received_validator_count() -> i64 {
    unsafe { raw::getReceivedValidatorCount() }
}

pub fn get_prepare_block_time() -> i64 {
    unsafe { raw::getPrepareBlockTime() }
}

pub fn get_aggregate_block_time() -> i64 {
    unsafe { raw::getAggregateBlockTime() }
}

pub fn get_validator_address(index: i64) -> Vec<u8> {
    unsafe {
        let mut data = vec![0u8; 20];
        assert_eq!(0, raw::readValidatorAddress(index, data.as_mut_ptr()));
        data
    }
}

pub fn get_calldata() -> Vec<u8> {
    unsafe {
        let data_size = raw::getCallDataSize();
        let mut data = vec![0u8; data_size as usize];
        assert_eq!(0, raw::readCallData(data.as_mut_ptr(), 0, data_size));
        data
    }
}

pub fn save_return_data(data: &[u8]) {
    unsafe { assert_eq!(0, raw::saveReturnData(data.as_ptr(), data.len() as i64)) }
}

pub fn request_external_data(data_source_id: i64, external_id: i64, calldata: &[u8]) {
    unsafe {
        assert_eq!(
            0,
            raw::requestExternalData(
                data_source_id,
                external_id,
                calldata.as_ptr(),
                calldata.len() as i64
            )
        )
    }
}

pub fn get_external_data(external_id: i64, validator_index: i64) -> Option<String> {
    unsafe {
        let data_size = raw::getExternalDataSize(external_id, validator_index);
        if data_size == -1 {
            None
        } else {
            let mut data = vec![0u8; data_size as usize];
            assert_eq!(
                0,
                raw::readExternalData(
                    external_id,
                    validator_index,
                    data.as_mut_ptr(),
                    0,
                    data_size
                )
            );
            Some(String::from_utf8_unchecked(data))
        }
    }
}
