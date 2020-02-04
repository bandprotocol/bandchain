use owasm::ext::finance::yahoo_finance;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub stock_symbol: String,
    }
}

decl_data! {
    pub struct Data {
        pub yahoo_finance: f32 = |params: &Parameter| yahoo_finance::Price::new(&params.stock_symbol),
    }
}

decl_result! {
    pub struct Result {
        pub stock_price_in_usd: u64,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    let mut total_stock_price = 0.0;
    for each in &data {
        total_stock_price += each.yahoo_finance;
    }
    let average_stock_price = total_stock_price / (data.len() as f32);
    Result { stock_price_in_usd: (average_stock_price * 100.0) as u64 }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter { stock_symbol: String::from("FB") };
        let data1 = Data { yahoo_finance: 100.0 };
        let data2 = Data { yahoo_finance: 200.0 };
        assert_eq!(execute(params, vec![data1, data2]), Result { stock_price_in_usd: 15000 });
    }

    #[test]
    fn test_call_real_price() {
        let params = Parameter { stock_symbol: String::from("FB") };
        let data = Data::build_from_local_env(&params).unwrap();
        println!("Current FB price (times 100) is {:?}", execute(params, vec![data]));
    }
}
