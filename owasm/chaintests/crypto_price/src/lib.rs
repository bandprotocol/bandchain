use owasm::oei;

fn parse_coingecko_symbol(symbol: &[u8]) -> &[u8] {
    let s = String::from_utf8(symbol.to_vec()).unwrap();
    (match s.as_str() {
        "BTC" => "bitcoin",
        "ETH" => "ethereum",
        _ => panic!("Unsupported coin!"),
    })
    .as_bytes()
}

fn parse_float(data: String) -> Option<f64> {
    data.parse::<f64>().ok()
}

#[no_mangle]
pub fn prepare() {
    let calldata = oei::get_calldata();
    // Coingecko data source
    oei::request_external_data(1, 1, parse_coingecko_symbol(&calldata));
    // Crypto compare source
    oei::request_external_data(2, 2, &calldata);
    // Binance source
    oei::request_external_data(3, 3, &calldata);
}

#[no_mangle]
pub fn execute() {
    let validator_count = oei::get_requested_validator_count();
    let mut sum: f64 = 0.0;
    let mut count: u64 = 0;
    for validator_index in 0..validator_count {
        let mut val = 0.0;
        let mut fail = false;
        for external_id in 1..4 {
            let data = oei::get_external_data(external_id, validator_index);
            if data.is_none() {
                fail = true;
                break;
            }
            let num = parse_float(data.unwrap());
            if num.is_none() {
                fail = true;
                break;
            }
            val += num.unwrap();
        }
        if !fail {
            sum += val / 3.0;
            count += 1;
        }
    }
    let result = (sum / (count as f64) * 100.0) as u64;
    oei::save_return_data(&result.to_be_bytes())
}
