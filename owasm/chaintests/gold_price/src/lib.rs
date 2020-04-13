use owasm::oei;

fn parse_float(data: String) -> Option<f64> {
    data.parse::<f64>().ok()
}

#[no_mangle]
pub fn prepare() {
    let calldata = oei::get_calldata();
    // Gold price data source
    oei::request_external_data(5, 1, &calldata);
}

#[no_mangle]
pub fn execute() {
    let validator_count = oei::get_requested_validator_count();
    let mut sum: f64 = 0.0;
    let mut count: u64 = 0;
    for validator_index in 0..validator_count {
        let mut val = 0.0;
        let mut fail = false;

        let data = oei::get_external_data(5, validator_index);
        if data.is_none() {
            fail = true;
        }
        let num = parse_float(data.unwrap());
        if num.is_none() {
            fail = true;
        }
        val += num.unwrap();

        if !fail {
            sum += val;
            count += 1;
        }
    }
    let result = (sum / (count as f64) * 100.0) as u64;
    oei::save_return_data(&result.to_be_bytes())
}
