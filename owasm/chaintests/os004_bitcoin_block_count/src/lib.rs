use obi::{get_schema, OBIDecode, OBIEncode, OBISchema};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode, OBISchema)]
struct Input {
    // TODO: remove this later
    _unused: u8,
}

#[derive(OBIEncode, OBISchema)]
struct Output {
    block_count: u64,
}

#[no_mangle]
fn prepare_impl(_: Input) {
    // Bitcoin block count data source
    oei::ask_external_data(1, 7, "".as_bytes());
}

#[no_mangle]
fn execute_impl(_: Input) -> Output {
    let majority = ext::load_majority::<u64>(1);
    Output { block_count: majority.unwrap() }
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
        assert_eq!("{_unused:u8}/{block_count:u64}", format!("{}/{}", input_schema, output_schema),);
    }
}
