use owasm::ext::bitcoin::{block_count, block_hash};
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub block_height: u64,
    }
}

decl_data! {
    pub struct Data {
        pub block_hash: String = |params: &Parameter| block_hash::Info::new(params.block_height),
        pub block_count: u64 = |_: &Parameter| block_count::Info::new(),
    }
}

decl_result! {
    pub struct Result {
        pub block_hash: String,
        pub confirmation: u64,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    for a in &data {
        let mut count = 0;
        for b in &data {
            if a.block_hash == b.block_hash && a.block_count == b.block_count {
                count = count + 1;
                if count > data.len() / 2 {
                    return Result {
                        block_hash: a.block_hash.to_string(),
                        confirmation: a.block_count - _params.block_height,
                    };
                }
            }
        }
    }

    Result { block_hash: "".into(), confirmation: 0 }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute_when_every_report_the_same_things() {
        let params = Parameter { block_height: 616047 };
        let data1 = Data {
            block_hash: String::from(
                "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b",
            ),
            block_count: 616058,
        };
        let data2 = Data {
            block_hash: String::from(
                "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b",
            ),
            block_count: 616058,
        };
        let data3 = Data {
            block_hash: String::from(
                "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b",
            ),
            block_count: 616058,
        };
        assert_eq!(
            execute(params, vec![data1, data2, data3]),
            Result {
                block_hash: String::from(
                    "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b"
                ),
                confirmation: 616058 - 616047
            }
        );
    }

    #[test]
    fn test_execute_when_the_minority_report_differently() {
        let params = Parameter { block_height: 616047 };
        let data1 = Data {
            block_hash: String::from(
                "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b",
            ),
            block_count: 616058,
        };
        let data2 = Data {
            block_hash: String::from(
                "0000000000000000000000000000000000000000000000000000000000000000",
            ),
            block_count: 616058,
        };
        let data3 = Data {
            block_hash: String::from(
                "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b",
            ),
            block_count: 616058,
        };
        assert_eq!(
            execute(params, vec![data1, data2, data3]),
            Result {
                block_hash: String::from(
                    "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b"
                ),
                confirmation: 616058 - 616047
            }
        );
    }

    #[test]
    fn test_execute_when_cant_find_majority() {
        let params = Parameter { block_height: 616047 };
        let data1 = Data {
            block_hash: String::from(
                "0000000000000000000c71e9f3636060e152b30fcc47c44b26a57f6c16c4447b",
            ),
            block_count: 616058,
        };
        let data2 = Data {
            block_hash: String::from(
                "0000000000000000000000000000000000000000000000000000000000000000",
            ),
            block_count: 616058,
        };
        assert_eq!(
            execute(params, vec![data1, data2]),
            Result { block_hash: String::from(""), confirmation: 0 }
        );
    }

    #[test]
    fn test_call_real_price() {
        let params = Parameter { block_height: 616047 };
        let data = Data::build_from_local_env(&params).unwrap();
        println!(
            "Bitcoin's block_hash and confirmation at block_height=616047 is {:?}",
            execute(params, vec![data])
        );
    }
}
