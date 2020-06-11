use obi::{get_schema, OBIDecode, OBIEncode, OBISchema};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode, OBISchema)]
struct Input {
    flight_op: String,
    airport: String,
    icao24: String,
    begin: String,
    end: String,
}

#[derive(OBIEncode, OBISchema)]
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

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::*;

    #[test]
    fn test_get_schema() {
        let mut schema = HashMap::new();
        Input::add_definitions_recursively(&mut schema);
        Output::add_definitions_recursively(&mut schema);
        let input_schema = get_schema(String::from("Input"), &schema);
        let output_schema = get_schema(String::from("Output"), &schema);
        println!("{}/{}", input_schema, output_schema);
        assert_eq!(
            "{flight_op:string,airport:string,icao24:string,begin:string,end:string}/{flight_existence:bool}",
            format!("{}/{}", input_schema, output_schema),
        );
    }
}
