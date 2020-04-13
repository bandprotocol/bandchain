use borsh::{BorshDeserialize, BorshSerialize};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(BorshDeserialize)]
struct Input {
    symbol: String,
    multiplier: u64,
}

#[derive(BorshSerialize)]
struct Output {
    px: u64,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Gold price data source
    oei::request_external_data(5, 1, input.symbol.as_bytes());
    // Binance data source
    oei::request_external_data(6, 2, input.symbol.as_bytes());
}

#[no_mangle]
fn execute_impl(input: Input) -> Output {
    let avg_gold: f64 = ext::load_average(1);
    let avg_atom: f64 = ext::load_average(2);
    Output { px: (avg_gold * input.multiplier as f64 / avg_atom) as u64 }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
