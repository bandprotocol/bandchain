use obi::{get_schema, OBIDecode, OBIEncode, OBISchema};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode, OBISchema)]
struct Input {
    gas_option: String,
}

#[derive(OBIEncode, OBISchema)]
struct Output {
    gweix10: u64,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // ETH gas station data source
    oei::ask_external_data(1, 11, &input.gas_option.as_bytes());
}

#[no_mangle]
fn execute_impl(_: Input) -> Output {
    let avg: f64 = ext::load_average(1);
    Output { gweix10: avg as u64 }
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
            "{gas_option:string}/{gweix10:u64}",
            format!("{}/{}", input_schema, output_schema),
        );
    }
}
