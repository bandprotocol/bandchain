mod raw;

/// Returns the number of validators to asked to report data from raw requests.
pub fn get_ask_count() -> i64 {
    unsafe { raw::get_ask_count() }
}

/// Returns the minimum number of data reports as specified by the oracle request.
pub fn get_min_count() -> i64 {
    unsafe { raw::get_min_count() }
}

/// Returns the number of validators that report data to this oracle request. Must
/// only be called during execution phase.
pub fn get_ans_count() -> i64 {
    unsafe { raw::get_ans_count() }
}

/// Returns the raw calldata as specified when the oracle request is submitted.
pub fn get_calldata() -> Vec<u8> {
    unsafe {
        let mut data = Vec::with_capacity(raw::get_span_size() as usize);
        let len = raw::read_calldata(data.as_mut_ptr() as i64);
        data.set_len(len as usize);
        data
    }
}

/// Saves the given data as the result of the oracle execution. Must only be called
/// during execution phase and must be called exactly once.
pub fn save_return_data(data: &[u8]) {
    unsafe { raw::set_return_data(data.as_ptr() as i64, data.len() as i64) }
}

/// Issues a new raw request to the host environement using the specified data
/// source ID and calldata, and assigns it to the given external ID. Must only be
/// called during preparation phase.
pub fn ask_external_data(eid: i64, did: i64, calldata: &[u8]) {
    unsafe { raw::ask_external_data(eid, did, calldata.as_ptr() as i64, calldata.len() as i64) }
}

/// Returns the data reported from the given validator index for the given external
/// data ID. Result is OK if the validator reports data with zero return status, and
/// Err otherwise. Must only be called during execution phase.
pub fn get_external_data(eid: i64, vid: i64) -> Result<String, i64> {
    unsafe {
        let status = raw::get_external_data_status(eid, vid);
        if status != 0 {
            Err(status)
        } else {
            let mut data = Vec::with_capacity(raw::get_span_size() as usize);
            let len = raw::read_external_data(eid, vid, data.as_mut_ptr() as i64);
            data.set_len(len as usize);
            Ok(String::from_utf8_unchecked(data))
        }
    }
}
