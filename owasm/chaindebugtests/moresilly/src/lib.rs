use obi::{get_schema, OBIDecode, OBIEncode, OBISchema};
use owasm::ext;
use owasm::oei;
use owasm::{execute_entry_point, prepare_entry_point};

#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

const PRICE_DATA_SOURCE_ID: i64 = 1;
const PRICE_EXTERNAL_ID: i64 = 1;

#[derive(OBIDecode, OBISchema)]
struct Input {
    symbol: String,
    multiplier: u64,
}

#[derive(OBIEncode, OBISchema)]
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
            "{symbol:string,multiplier:u64}/{px:u64}",
            format!("{}/{}", input_schema, output_schema),
        );
    }
}
