use owasm::ext::ethgasstation::gas_price;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub gas_option: String,
    }
}

decl_data! {
    pub struct Data {
        pub gas_price: f32 = |params: &Parameter| gas_price::Price::new(&params.gas_option),
    }
}

decl_result! {
    pub struct Result {
        pub gas_price_in_gwei: u64,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    let mut total_gas_price = 0.0;
    for each in &data {
        total_gas_price += each.gas_price;
    }
    let average_gas_price = total_gas_price / (data.len() as f32);
    Result { gas_price_in_gwei: (average_gas_price * 100.0) as u64 }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter { gas_option: String::from("average") };
        let data1 = Data { gas_price: 13.0 };
        let data2 = Data { gas_price: 7.0 };
        assert_eq!(execute(params, vec![data1, data2]), Result { gas_price_in_gwei: 1000 });
    }

    #[test]
    fn test_call_real_gas_price() {
        let params = Parameter { gas_option: String::from("average") };
        let data = Data::build_from_local_env(&params).unwrap();
        println!(
            "Current Ethereum gas price with average option (times 100) is {:?}",
            execute(params, vec![data])
        );
    }
}
