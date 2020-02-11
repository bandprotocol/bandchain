use owasm::ext::crypto::{binance, coins};
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub crypto_symbol: coins::Coins,
    }
}

decl_data! {
    pub struct Data {
        pub binance_price: f32 = |params: &Parameter| binance::Price::new(&params.crypto_symbol),
        pub timestamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub crypto_price_in_usd: u64,
        pub timestamp: u64,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    let mut total_crypto_price = 0.0;
    let mut timestamp_acc: u64 = 0;
    for each in &data {
        total_crypto_price += each.binance_price;
        timestamp_acc += each.timestamp;
    }
    let average_crypto_price = total_crypto_price / (data.len() as f32);
    let avg_timestamp = timestamp_acc / (data.len() as u64);
    Result { crypto_price_in_usd: (average_crypto_price * 100.0) as u64, timestamp: avg_timestamp }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter { crypto_symbol: coins::Coins::ETH };
        let data1 = Data { binance_price: 100.0, timestamp: 10 };
        let data2 = Data { binance_price: 200.0, timestamp: 12 };
        assert_eq!(
            execute(params, vec![data1, data2]),
            Result { crypto_price_in_usd: 15000, timestamp: 11 }
        );
    }

    #[test]
    fn test_call_real_price() {
        let params = Parameter { crypto_symbol: coins::Coins::ETH };
        let data = Data::build_from_local_env(&params).unwrap();
        println!("Current ETH price (times 100) is {:?}", execute(params, vec![data]));
    }
}
