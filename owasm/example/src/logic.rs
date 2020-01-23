use owasm::ext::crypto::{binance, coingecko, coins, cryptocompare};
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub symbol: coins::Coins,
    }
}

decl_data! {
    pub struct Data {
        pub coin_gecko: f32 = |params: &Parameter| coingecko::Price::new(&params.symbol),
        pub crypto_compare: f32 = |params: &Parameter| cryptocompare::Price::new(&params.symbol),
        pub binance: f32 = |params: &Parameter| binance::Price::new(&params.symbol),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub price_in_usd: u64,
        pub time_stamp: u64,
    }
}

impl Data {
    pub fn avg_px(&self) -> f32 {
        (self.coin_gecko + self.crypto_compare + self.binance) / 3.0
    }
}

pub fn execute(data: Vec<Data>) -> Result {
    let mut total = 0.0;
    let mut time_stamp_acc: u64 = 0;
    for each in &data {
        total += each.avg_px();
        time_stamp_acc += each.time_stamp;
    }
    let average = total / (data.len() as f32);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result { price_in_usd: (average * 100.0) as u64, time_stamp: avg_time_stamp }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        // Average is 120.00
        let data1 =
            Data { coin_gecko: 100.0, crypto_compare: 150.0, binance: 110.0, time_stamp: 10 };
        // Average is 220.00
        let data2 =
            Data { coin_gecko: 200.0, crypto_compare: 250.0, binance: 210.0, time_stamp: 12 };
        // Average among the two data points is 170.00
        assert_eq!(execute(vec![data1, data2]), Result { price_in_usd: 17000, time_stamp: 11 });
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
