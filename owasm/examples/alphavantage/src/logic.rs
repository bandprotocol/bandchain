use owasm::ext::finance::alphavantage;
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub stock_symbol: String,
    }
}

decl_data! {
    pub struct Data {
        pub alphavantage_price: f32 = |params: &Parameter| alphavantage::Price::new(&params.stock_symbol, "WVKPOO76169EX950"),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub stock_price_in_usd: u64,
        pub time_stamp: u64,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    let mut total_stock_price = 0.0;
    let mut time_stamp_acc: u64 = 0;
    for each in &data {
        total_stock_price += each.alphavantage_price;
        time_stamp_acc += each.time_stamp;
    }
    let average_stock_price = total_stock_price / (data.len() as f32);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result { stock_price_in_usd: (average_stock_price * 100.0) as u64, time_stamp: avg_time_stamp }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter { stock_symbol: String::from("FB") };
        let data1 = Data { alphavantage_price: 100.0, time_stamp: 10 };
        let data2 = Data { alphavantage_price: 200.0, time_stamp: 12 };
        assert_eq!(
            execute(params, vec![data1, data2]),
            Result { stock_price_in_usd: 15000, time_stamp: 11 }
        );
    }

    #[test]
    fn test_call_real_price() {
        let params = Parameter { stock_symbol: String::from("FB") };
        let data = Data::build_from_local_env(&params).unwrap();
        println!("Current FB price (times 100) is {:?}", execute(params, vec![data]));
    }
}
