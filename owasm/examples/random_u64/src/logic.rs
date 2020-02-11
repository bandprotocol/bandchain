use owasm::ext::random::qrng_anu;
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

use std::convert::TryInto;

decl_params! {
    pub struct Parameter {
        pub max_range: u64,
    }
}

decl_data! {
    pub struct Data {
        pub random_bytes8: Vec<u8> = |_: &Parameter| qrng_anu::RandomBytes::new(8),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub random_u64: u64,
        pub time_stamp: u64,
    }
}

impl Data {
    pub fn rng_to_u64(&self) -> u64 {
        u64::from_le_bytes(
            TryInto::<[u8; 8]>::try_into((&self.random_bytes8) as &[u8]).unwrap_or([0; 8]),
        )
    }
}

pub fn execute(params: Parameter, data: Vec<Data>) -> Result {
    Result {
        random_u64: data.iter().fold(0, |accumulator, each| accumulator ^ each.rng_to_u64())
            % params.max_range,
        time_stamp: data.iter().fold(0, |accumulator, each| accumulator + each.time_stamp)
            / (data.len() as u64),
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter { max_range: 31 };
        let data1 = Data { random_bytes8: vec![1, 2, 3, 4, 5, 6, 7, 8], time_stamp: 10 };
        let data2 = Data { random_bytes8: vec![8, 7, 6, 5, 4, 3, 2, 1], time_stamp: 12 };
        assert_eq!(execute(params, vec![data1, data2]), Result { random_u64: 18, time_stamp: 11 });
    }

    #[test]
    fn test_call_get_random_u64() {
        let params = Parameter { max_range: 31 };
        let data = Data::build_from_local_env(&params)
            .unwrap_or(Data { random_bytes8: vec![0; 8], time_stamp: 0 });
        println!("Current random number is {:?}", execute(params, vec![data]));
    }
}
