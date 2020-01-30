use owasm::ext::random::qrng_anu;
use owasm::{decl_data, decl_params, decl_result};

use std::convert::TryInto;

decl_params! {
    pub struct Parameter {}
}

decl_data! {
    pub struct Data {
        pub rng: Vec<u8> = |_: &Parameter| qrng_anu::RandomBytes::new(8),
    }
}

decl_result! {
    pub struct Result {
        pub random_number: u64,
    }
}

impl Data {
    pub fn rng_to_u64(&self) -> u64 {
        match ((&self.rng) as &[u8]).try_into().ok() {
            Some(data) => u64::from_le_bytes(data),
            None => 0,
        }
    }
}

pub fn execute(data: Vec<Data>) -> Result {
    let mut acc_rng = 0;
    for each in &data {
        acc_rng ^= each.rng_to_u64();
    }
    Result { random_number: acc_rng }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let data1 = Data { rng: vec![1, 2, 3, 4, 5, 6, 7, 8] };
        let data2 = Data { rng: vec![8, 7, 6, 5, 4, 3, 2, 1] };
        assert_eq!(execute(vec![data1, data2]), Result { random_number: 649931223095117065 });
    }
}
