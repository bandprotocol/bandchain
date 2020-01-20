use owasm::ext::crypto::{coingecko, coins, cryptocompare};
use owasm::{decl_data, decl_params, decl_result};
use serde::{Deserialize, Serialize};

decl_params! {
    pub struct Parameter {
        pub symbol: coins::Coins,
    }
}

decl_data! {
    pub struct Data {
        pub coin_gecko: f32 = |params: &Parameter| coingecko::Price::new(&params.symbol),
        pub crypto_compare: f32 = |params: &Parameter| cryptocompare::Price::new(&params.symbol),
    }
}

decl_result! {
    pub struct Result {
        pub price_in_usd: u64,
    }
}

impl Data {
    pub fn avg_px(&self) -> f32 {
        (self.coin_gecko + self.crypto_compare) / 2.0
    }
}

pub fn execute(data: Vec<Data>) -> Result {
    let mut total = 0.0;
    for each in &data {
        total += each.avg_px();
    }
    let average = total / (data.len() as f32);
    Result { price_in_usd: (average * 100.0) as u64 }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        // Average is 125.00
        let data1 = Data { coin_gecko: 100.0, crypto_compare: 150.0 };
        // Average is 225.00
        let data2 = Data { coin_gecko: 200.0, crypto_compare: 250.0 };
        // Average among the two data points is 175.00
        assert_eq!(execute(vec![data1, data2]), Result { price_in_usd: 17500 });
    }

    #[test]
    fn test_end_to_end_from_local_env() {
        // Run with local environment
        let data = Data::build_from_local_env(&Parameter { symbol: coins::Coins::BTC }).unwrap();
        println!("Current BTC price (times 100) is {:?}", execute(vec![data]));

        let data = Data::build_from_local_env(&Parameter { symbol: coins::Coins::ETH }).unwrap();
        println!("Current ETH price (times 100) is {:?}", execute(vec![data]));
    }
}
