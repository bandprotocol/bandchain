use obi::{OBIDecode, OBIEncode};
use owasm::{execute_entry_point, oei, prepare_entry_point};

fn parse_float(data: String) -> Option<f64> {
    data.trim_end().parse::<f64>().ok()
}

#[derive(OBIDecode)]
struct Input {
    symbol: String,
    multiplier: u64,
}

#[derive(OBIEncode)]
struct Output {
    px: u64,
}

fn prepare_impl(input: Input) {
    // Coingecko data source
    oei::request_external_data(1, 1, &input.symbol.as_bytes());
    // Crypto compare source
    oei::request_external_data(2, 2, &input.symbol.as_bytes());
    // Binance source
    oei::request_external_data(3, 3, &input.symbol.as_bytes());
}

fn execute_impl(input: Input) -> Output {
    let validator_count = oei::get_ask_count();
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
    Output { px: (sum / (count as f64) * (input.multiplier as f64)) as u64 }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_get_schema() {
        let input = Input { symbol: String::from("BTC"), multiplier: 100 };
        let encoded_calldata: [u8; 15] = [0, 0, 0, 3, 66, 84, 67, 0, 0, 0, 0, 0, 0, 0, 100];
        let result: Input = OBIDecode::try_from_slice(&encoded_calldata).unwrap();
        assert_eq!(input.multiplier, result.multiplier);
        assert_eq!(input.symbol, result.symbol);
    }
}
