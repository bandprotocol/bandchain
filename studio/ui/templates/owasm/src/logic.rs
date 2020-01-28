use owasm::ext::crypto::{binance, coingecko, coins, cryptocompare};
use owasm::ext::finance::{alphavantage};
use owasm::ext::random::qrng_anu;
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

use std::convert::TryInto;

decl_params! {
    pub struct Parameter {
        pub crypto_symbol: coins::Coins,
        pub stock_symbol: String,
        pub alphavantage_api_key: String,
    }
}

decl_data! {
    pub struct Data {
        pub coin_gecko: f32 = |params: &Parameter| coingecko::Price::new(&params.crypto_symbol),
        pub coin_gecko_vol24h: f64 = |params: &Parameter| coingecko::Volume24h::new(&params.crypto_symbol),
        pub crypto_compare: f32 = |params: &Parameter| cryptocompare::Price::new(&params.crypto_symbol),
        pub crypto_compare_vol24h: f64 = |params: &Parameter| cryptocompare::Volume24h::new(&params.crypto_symbol),
        pub binance: f32 = |params: &Parameter| binance::Price::new(&params.crypto_symbol),
        pub alphavantage: f32 = |params: &Parameter| alphavantage::Price::new(&params.stock_symbol, &params.alphavantage_api_key),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
        pub rng: Vec<u8> = |_: &Parameter| qrng_anu::RandomBytes::new(8),
    }
}

decl_result! {
    pub struct Result {
        pub crypto_price_in_usd: u64,
        pub time_stamp: u64,
        pub stock_price_in_usd: u64,
        pub random_number: u64,
        pub vol24h_in_usd: u64,
    }
}

impl Data {
    pub fn avg_vol24h(&self) -> f64 {
        (self.coin_gecko_vol24h + self.crypto_compare_vol24h) / 2.0
    }

    pub fn avg_crypto_px(&self) -> f32 {
        (self.coin_gecko + self.crypto_compare + self.binance) / 3.0
    }

    pub fn avg_stock_px(&self) -> f32 {
        self.alphavantage
    }

    pub fn rng_to_u64(&self) -> u64 {
        match ((&self.rng) as &[u8]).try_into().ok() {
            Some(data) => u64::from_le_bytes(data),
            None => 0,
        }
    }
}

pub fn execute(data: Vec<Data>) -> Result {
    let mut total_crypto_price = 0.0;
    let mut total_stock_price = 0.0;
    let mut time_stamp_acc: u64 = 0;
    let mut acc_rng = 0;
    let mut total_vol: f64 = 0.0;
    for each in &data {
        total_crypto_price += each.avg_crypto_px();
        total_stock_price += each.avg_stock_px();
        time_stamp_acc += each.time_stamp;
        total_vol += each.avg_vol24h();
        acc_rng ^= each.rng_to_u64();
    }
    let average_crypto_price = total_crypto_price / (data.len() as f32);
    let average_stock_price = total_stock_price / (data.len() as f32);
    let average_vol = total_vol / (data.len() as f64);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result {
        crypto_price_in_usd: (average_crypto_price * 100.0) as u64,
        stock_price_in_usd: (average_stock_price * 100.0) as u64,
        vol24h_in_usd: (average_vol * 100.0) as u64,
        time_stamp: avg_time_stamp,
        random_number: acc_rng,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        // Average of crypto price is 120.00
        // Average of stock price is 220.00
        let data1 = Data {
            coin_gecko: 100.0,
            crypto_compare: 150.0,
            binance: 110.0,
            alphavantage: 230.0,
            coin_gecko_vol24h: 200.0,
            crypto_compare_vol24h: 250.0,
            time_stamp: 10,
            rng: vec![1, 2, 3, 4, 5, 6, 7, 8],
        };
        // Average of crypto price is 220.00
        // Average of stock price is 320.00
        let data2 = Data {
            coin_gecko: 200.0,
            crypto_compare: 250.0,
            binance: 210.0,
            alphavantage: 310.0,
            coin_gecko_vol24h: 150.0,
            crypto_compare_vol24h: 100.0,
            time_stamp: 12,
            rng: vec![8, 7, 6, 5, 4, 3, 2, 1],
        };
        // Average among the two crypto data points is 170.00
        // Average among the two stock data points is 270.00
        assert_eq!(
            execute(vec![data1, data2]),
            Result {
                crypto_price_in_usd: 17000,
                stock_price_in_usd: 27000,
                vol24h_in_usd: 17500,
                time_stamp: 11,
                random_number: 649931223095117065
            }
        );
    }

    #[test]
    fn test_call_stock() {
        let data = Data::build_from_local_env(&Parameter {
            crypto_symbol: coins::Coins::ETH,
            stock_symbol: String::from("FB"),
            alphavantage_api_key: String::from("WVKPOO76169EX950"),
        })
        .unwrap();
        println!("Current ETH and FB price (times 100) is {:?}", execute(vec![data]));
    }
}