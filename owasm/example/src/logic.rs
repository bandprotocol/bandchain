use owasm::ext::crypto::{alphavantage, binance, coingecko, coins, cryptocompare};
use owasm::ext::finance::alphavantage;
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub symbol: coins::Coins,
        pub alphavantage_symbol: String,
        pub alphavantage_api_key: String,
    }
}

decl_data! {
    pub struct Data {
        pub coin_gecko: f32 = |params: &Parameter| coingecko::Price::new(&params.symbol),
        pub crypto_compare: f32 = |params: &Parameter| cryptocompare::Price::new(&params.symbol),
        pub binance: f32 = |params: &Parameter| binance::Price::new(&params.symbol),
        pub alphavantage: f32 = |params: &Parameter| alphavantage::Price::new(&params.alphavantage_symbol, &params.alphavantage_api_key),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub price_in_usd: u64,
        pub time_stamp: u64,
        pub price_from_alphavantage_in_usd: u64,
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
    let mut total_price_from_alphavantage = 0.0;
    for each in &data {
        total += each.avg_px();
        time_stamp_acc += each.time_stamp;
        total_price_from_alphavantage += each.alphavantage;
    }
    let average = total / (data.len() as f32);
    let average_price_from_alphavantage = total_price_from_alphavantage / (data.len() as f32);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result {
        price_in_usd: (average * 100.0) as u64,
        price_from_alphavantage_in_usd: (average_price_from_alphavantage * 100.0) as u64,
        time_stamp: avg_time_stamp,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        // Average is 120.00
        let data1 = Data {
            coin_gecko: 100.0,
            crypto_compare: 150.0,
            binance: 110.0,
            alphavantage: 120.0,
            time_stamp: 10,
        };
        // Average is 220.00
        let data2 = Data {
            coin_gecko: 200.0,
            crypto_compare: 250.0,
            binance: 210.0,
            alphavantage: 220.0,
            time_stamp: 12,
        };
        // Average among the two data points is 170.00
        assert_eq!(
            execute(vec![data1, data2]),
            Result { price_in_usd: 17000, price_from_alphavantage_in_usd: 17000, time_stamp: 11 }
        );
    }

    #[test]
    fn test_call_stock() {
        let data = Data::build_from_local_env(&Parameter {
            symbol: coins::Coins::ETH,
            alphavantage_symbol: String::from("GOOG"),
            alphavantage_api_key: String::from("WVKPOO76169EX950"),
        })
        .unwrap();
        println!("Current ETH and GOOG price (times 100) is {:?}", execute(vec![data]));
    }
}
