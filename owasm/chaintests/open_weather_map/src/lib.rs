use obi::{OBIDecode, OBIEncode};
use owasm2::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode)]
struct Input {
    country: String,
    main_field: String,
    sub_field: String,
    multiplier: u64,
}

#[derive(OBIEncode)]
struct Output {
    value: u64,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Open weather data source
    let Input { country, main_field, sub_field, .. } = input;
    oei::ask_external_data(1, 4, format!("{} {} {}", country, main_field, sub_field).as_bytes());
}

#[no_mangle]
fn execute_impl(input: Input) -> Output {
    let avg: f64 = ext::load_average(1);
    Output { value: (avg * input.multiplier as f64) as u64 }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
