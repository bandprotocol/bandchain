use obi::{OBIDecode, OBIEncode};
use owasm2::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode)]
struct Input {
    flight_op: String,
    airport: String,
    icao24: String,
    begin: String,
    end: String,
}

#[derive(OBIEncode)]
struct Output {
    flight_existence: bool,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Open sky api data source
    let Input { flight_op, airport, icao24, begin, end } = input;
    oei::ask_external_data(
        1,
        12,
        format!("{} {} {} {} {}", flight_op, airport, icao24, begin, end).as_bytes(),
    );
}

#[no_mangle]
fn execute_impl(_: Input) -> Output {
    let major = ext::load_majority::<bool>(1);
    Output { flight_existence: major.unwrap() }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
