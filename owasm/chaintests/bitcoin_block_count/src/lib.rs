use borsh::{BorshDeserialize, BorshSchema, BorshSerialize};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(BorshDeserialize, BorshSchema)]
struct Input {
    // TODO: remove this later
    _unused: u8,
}

#[derive(BorshSerialize, BorshSchema)]
struct Output {
    block_count: u64,
}

#[no_mangle]
fn prepare_impl(_: Input) {
    // Bitcoin block count data source
    oei::request_external_data(6, 1, "".as_bytes());
}

#[no_mangle]
fn execute_impl(_: Input) -> Output {
    let majority = ext::load_majority::<u64>(1);
    Output { block_count: majority.unwrap() }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
