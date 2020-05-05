use borsh::{BorshDeserialize, BorshSchema, BorshSerialize};
use owasm::{execute_entry_point, oei};

#[derive(BorshDeserialize, BorshSchema)]
struct Input {
    _unused: u8,
}

#[derive(BorshSerialize, BorshSchema)]
struct Output {
    result: String,
}

// Expect fail when prepare
#[no_mangle]
fn prepare(_input: Input) {
    oei::request_external_data(1, 1, "Hello world".as_bytes());
}

fn execute_impl(_input: Input) -> Output {
    Output { result: String::from("Yeah") }
}

// prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::*;

    #[test]
    fn test_get_schema() {
        let mut schema = HashMap::new();
        Input::add_rec_type_definitions(&mut schema);
        Output::add_rec_type_definitions(&mut schema);
        println!("{:?}", schema);
    }
}
