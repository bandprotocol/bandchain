use owasm::ext::crypto::{binance, coingecko, coins, cryptocompare};
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub crypto_symbol: coins::Coins,
    }
}

decl_data! {
    pub struct Data {
        pub eth_coin_gecko: f32 = |_: &Parameter| coingecko::Price::new(&coins::Coins::ETH),
        pub eth_crypto_compare: f32 = |_: &Parameter| cryptocompare::Price::new(&coins::Coins::ETH),
        pub eth_binance: f32 = |_: &Parameter| binance::Price::new(&coins::Coins::ETH),
        pub other_coin_gecko: f32 = |params: &Parameter| coingecko::Price::new(&params.crypto_symbol),
        pub other_crypto_compare: f32 = |params: &Parameter| cryptocompare::Price::new(&params.crypto_symbol),
        pub other_binance: f32 = |params: &Parameter| binance::Price::new(&params.crypto_symbol),
    }
}

decl_result! {
    pub struct Result {
        pub eth_price_in_usd: u64,
        pub other_price_in_usd: u64,
    }
}

impl Data {
    pub fn avg_crypto_px(&self) -> (f32, f32) {
        ((self.eth_coin_gecko + self.eth_crypto_compare + self.eth_binance) / 3.0,
        (self.other_coin_gecko + self.other_crypto_compare + self.other_binance) / 3.0)
    }
}

pub fn execute(data: Vec<Data>) -> Result {
    let mut total_eth_price = 0.0;
    let mut total_other_price = 0.0;
    for each in &data {
        let (eth_price, other_price) =  each.avg_crypto_px();
        total_eth_price += eth_price;
        total_other_price += other_price;
    }
    let average_eth_price = total_eth_price / (data.len() as f32);
    let average_other_price = total_other_price / (data.len() as f32);
    Result {
        eth_price_in_usd: (average_eth_price * 100.0) as u64,
        other_price_in_usd: (average_other_price * 100.0) as u64,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let data1 = Data {
            eth_coin_gecko: 100.0,
            eth_crypto_compare: 150.0,
            eth_binance: 110.0,
            other_coin_gecko: 10.0,
            other_crypto_compare: 15.0,
            other_binance: 20.0,
        };
        let data2 = Data {
            eth_coin_gecko: 200.0,
            eth_crypto_compare: 250.0,
            eth_binance: 210.0,
            other_coin_gecko: 12.0,
            other_crypto_compare: 13.0,
            other_binance: 14.0,
        };
        assert_eq!(
            execute(vec![data1, data2]),
            Result {
                eth_price_in_usd: 17000,
                other_price_in_usd: 1400,
            }
        );
    }

    #[test]
    fn test_price() {
        let data = Data::build_from_local_env(&Parameter {
            crypto_symbol: coins::Coins::EOS,
        })
        .unwrap();
        println!("Current ETH and EOS price (times 100) is {:?}", execute(vec![data]));
    }
}
