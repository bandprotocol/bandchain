use obi::{OBIDecode, OBIEncode};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode)]
struct Input {
    block_height: u64,
}

#[derive(OBIEncode)]
struct Output {
    block_hash: String,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Bitcoin hash data source
    oei::request_external_data(8, 1, input.block_height.to_string().as_bytes());
}

#[no_mangle]
fn execute_impl(_: Input) -> Output {
    let majority = ext::load_majority::<String>(1);
    Output { block_hash: majority.unwrap() }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
