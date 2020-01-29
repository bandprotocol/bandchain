use owasm::ext::crypto::{coingecko, coins, cryptocompare};
use owasm::ext::finance::{alphavantage};
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

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
        pub crypto_compare: f32 = |params: &Parameter| cryptocompare::Price::new(&params.crypto_symbol),
        pub alphavantage: f32 = |params: &Parameter| alphavantage::Price::new(&params.stock_symbol, &params.alphavantage_api_key),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub crypto_price_in_usd: u64,
        pub time_stamp: u64,
        pub stock_price_in_usd: u64,
    }
}

impl Data {
    pub fn avg_crypto_px(&self) -> f32 {
        (self.coin_gecko + self.crypto_compare) / 2.0
    }

    pub fn avg_stock_px(&self) -> f32 {
        self.alphavantage
    }
}

pub fn execute(data: Vec<Data>) -> Result {
    let mut total_crypto_price = 0.0;
    let mut total_stock_price = 0.0;
    let mut time_stamp_acc: u64 = 0;
    // let mut acc_rng = 0;
    for each in &data {
        total_crypto_price += each.avg_crypto_px();
        total_stock_price += each.avg_stock_px();
        time_stamp_acc += each.time_stamp;
    }
    let average_crypto_price = total_crypto_price / (data.len() as f32);
    let average_stock_price = total_stock_price / (data.len() as f32);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result {
        crypto_price_in_usd: (average_crypto_price * 100.0) as u64,
        stock_price_in_usd: (average_stock_price * 100.0) as u64,
        time_stamp: avg_time_stamp,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let data1 = Data {
            coin_gecko: 100.0,
            crypto_compare: 150.0,
            alphavantage: 230.0,
            time_stamp: 10,
        };
        let data2 = Data {
            coin_gecko: 200.0,
            crypto_compare: 250.0,
            alphavantage: 310.0,
            time_stamp: 12,
        };
        assert_eq!(
            execute(vec![data1, data2]),
            Result {
                crypto_price_in_usd: 17500,
                stock_price_in_usd: 27000,
                time_stamp: 11,
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
