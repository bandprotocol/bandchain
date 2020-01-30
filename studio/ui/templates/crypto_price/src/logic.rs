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
    pub fn avg_crypto_px(&self) -> f32 {
        (self.coin_gecko + self.crypto_compare + self.binance) / 3.0
    }
}

pub fn execute(data: Vec<Data>) -> Result {
    let mut total_crypto_price = 0.0;
    let mut time_stamp_acc: u64 = 0;
    for each in &data {
        total_crypto_price += each.avg_crypto_px();
        time_stamp_acc += each.time_stamp;
    }
    let average_crypto_price = total_crypto_price / (data.len() as f32);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result { price_in_usd: (average_crypto_price * 100.0) as u64, time_stamp: avg_time_stamp }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let data1 =
            Data { coin_gecko: 100.0, crypto_compare: 150.0, binance: 230.0, time_stamp: 10 };
        let data2 =
            Data { coin_gecko: 200.0, crypto_compare: 250.0, binance: 310.0, time_stamp: 12 };
        assert_eq!(execute(vec![data1, data2]), Result { price_in_usd: 20666, time_stamp: 11 });
    }

    #[test]
    fn test_call_eth_price() {
        let data = Data::build_from_local_env(&Parameter { symbol: coins::Coins::ETH }).unwrap();
        println!("Current ETH price (times 100) is {:?}", execute(vec![data]));
    }
}
