use obi::{get_schema, OBIDecode, OBIEncode, OBISchema};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode, OBISchema)]
struct Input {
    block_height: u64,
}

#[derive(OBIEncode, OBISchema)]
struct Output {
    block_hash: String,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Bitcoin hash data source
    oei::ask_external_data(1, 8, input.block_height.to_string().as_bytes());
}

#[no_mangle]
fn execute_impl(_: Input) -> Output {
    let majority = ext::load_majority::<String>(1);
    Output { block_hash: majority.unwrap() }
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
            "{block_height:u64}/{block_hash:string}",
            format!("{}/{}", input_schema, output_schema),
        );
    }
}
