use obi::{get_schema, OBIDecode, OBIEncode, OBISchema};
use owasm::{execute_entry_point, ext, oei};

#[derive(OBIDecode, OBISchema)]
struct Input {
    multiplier: u64,
}

#[derive(OBIEncode, OBISchema)]
struct Output {
    px: u64,
}

#[no_mangle]
fn prepare() {
    // Gold price data source
    oei::ask_external_data(1, 5, "".as_bytes());
    // Binance data source
    oei::ask_external_data(2, 3, "ATOM".as_bytes());
}

#[no_mangle]
fn execute_impl(input: Input) -> Output {
    let avg_gold: f64 = ext::load_average(1);
    let avg_atom: f64 = ext::load_average(2);
    Output { px: (avg_gold * input.multiplier as f64 / avg_atom) as u64 }
}

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
        assert_eq!("{multiplier:u64}/{px:u64}", format!("{}/{}", input_schema, output_schema),);
    }
}
