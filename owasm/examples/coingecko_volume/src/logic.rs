use owasm::ext::crypto::{coingecko, coins};
use owasm::ext::utils::date;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub crypto_symbol: coins::Coins,
    }
}

decl_data! {
    pub struct Data {
        pub coin_gecko_vol24h: f64 = |params: &Parameter| coingecko::Volume24h::new(&params.crypto_symbol),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub vol24h_in_usd: u64,
        pub time_stamp: u64,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    let mut total_vol: f64 = 0.0;
    let mut time_stamp_acc: u64 = 0;
    for each in &data {
        total_vol += each.coin_gecko_vol24h;
        time_stamp_acc += each.time_stamp;
    }
    let average_vol = total_vol / (data.len() as f64);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result { vol24h_in_usd: (average_vol * 100.0) as u64, time_stamp: avg_time_stamp }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter { crypto_symbol: coins::Coins::ETH };
        let data1 = Data { coin_gecko_vol24h: 250.0, time_stamp: 10 };
        let data2 = Data { coin_gecko_vol24h: 600.0, time_stamp: 12 };
        assert_eq!(
            execute(params, vec![data1, data2]),
            Result { vol24h_in_usd: 42500, time_stamp: 11 }
        );
    }

    #[test]
    fn test_call_real_volume() {
        let params = Parameter { crypto_symbol: coins::Coins::ETH };
        let data = Data::build_from_local_env(&params).unwrap();
        println!("Current ETH volume (times 100) is {:?}", execute(params, vec![data]));
    }
}
