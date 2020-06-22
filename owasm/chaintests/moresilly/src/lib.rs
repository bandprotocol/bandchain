use obi::{OBIDecode, OBIEncode};
use owasm2::ext;
use owasm2::oei;
use owasm2::{execute_entry_point, prepare_entry_point};

#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

const PRICE_DATA_SOURCE_ID: i64 = 1;
const PRICE_EXTERNAL_ID: i64 = 1;

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
    oei::ask_external_data(PRICE_EXTERNAL_ID, PRICE_DATA_SOURCE_ID, input.symbol.as_bytes());
}

fn execute_impl(input: Input) -> Output {
    let average_px: f64 = ext::load_average(PRICE_EXTERNAL_ID);
    Output { px: (average_px * input.multiplier as f64) as u64 }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
