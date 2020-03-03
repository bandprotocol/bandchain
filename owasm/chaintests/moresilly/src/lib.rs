use borsh::{BorshDeserialize, BorshSerialize};
use owasm::ext;
use owasm::oei;
use owasm::{execute_entry_point, prepare_entry_point};

#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

const PRICE_DATA_SOURCE_ID: i64 = 1;
const PRICE_EXTERNAL_ID: i64 = 1;

#[derive(BorshDeserialize)]
struct Input {
    symbol: String,
    multiplier: u64,
}

#[derive(BorshSerialize)]
struct Output {
    px: u64,
}

fn prepare_impl(input: Input) {
    oei::request_external_data(PRICE_DATA_SOURCE_ID, PRICE_EXTERNAL_ID, input.symbol.as_bytes());
}

fn execute_impl(input: Input) -> Output {
    let average_px: f64 = ext::load_average(PRICE_EXTERNAL_ID);
    Output { px: (average_px * input.multiplier as f64) as u64 }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
