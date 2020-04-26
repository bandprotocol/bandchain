use borsh::{BorshDeserialize, BorshSchema, BorshSerialize};
use hex;
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(BorshDeserialize, BorshSchema)]
struct Input {
    size: u64,
}

#[derive(BorshSerialize, BorshSchema)]
struct Output {
    random_bytes: String,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Quantum random data source
    oei::request_external_data(13, 1, input.size.to_string().as_bytes());
}

fn accumulate_hex_strings(strings: Vec<String>, input_size: usize) -> String {
    hex::encode(
        strings
            .iter()
            .map(|x1| x1.split(",").map(|x2| x2.parse::<u8>().unwrap()).collect::<Vec<_>>())
            .fold(vec![0; input_size], |acc, x1| {
                acc.iter().zip(x1.iter()).map(|x2| x2.0 ^ x2.1).collect::<Vec<_>>()
            }),
    )
}

#[no_mangle]
fn execute_impl(input: Input) -> Output {
    Output {
        random_bytes: accumulate_hex_strings(ext::load_input::<String>(1), input.size as usize),
    }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
