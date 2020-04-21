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

#[no_mangle]
fn execute_impl(_: Input) -> Output {
    let outputs = (ext::load_input::<String>(1))
        .iter()
        .map(|x1| x1.split(",").map(|x2| x2.parse::<u8>().unwrap()).collect::<Vec<_>>())
        .collect::<Vec<_>>();

    let mut acc = outputs[0].clone();
    for i in 1..(outputs.len()) {
        acc = acc.iter().zip(outputs[i].iter()).map(|x| x.0 ^ x.1).collect::<Vec<_>>();
    }

    Output { random_bytes: hex::encode(acc) }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);
