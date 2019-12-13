use owasm::decl_data;
use owasm::ext::{coingecko, cryptocompare};

decl_data! {
    pub struct Data {
        pub coin_gecko: f32 = coingecko::Price::new("bitcoin"),
        pub crypto_compare: f32 = cryptocompare::Price::new(),
    }
}

impl Data {
    pub fn avg_px(&self) -> f32 {
        (self.coin_gecko + self.crypto_compare) / 2.0
    }
}

pub fn execute(data: Vec<Data>) -> u64 {
    let mut total = 0.0;
    for each in &data {
        total += each.avg_px();
    }
    let average = total / (data.len() as f32);
    (average * 100.0) as u64
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        // Average is 125.00
        let data1 = Data {
            coin_gecko: 100.0,
            crypto_compare: 150.0,
        };
        // Average is 225.00
        let data2 = Data {
            coin_gecko: 200.0,
            crypto_compare: 250.0,
        };
        // Average among the two data points is 175.00
        assert_eq!(execute(vec![data1, data2]), 17500);
    }

    #[test]
    fn test_end_to_end_from_local_env() {
        // Run with local environment
        let data = Data::build_from_local_env().unwrap();
        println!("Current BTC price (times 100) is {}", execute(vec![data]));
    }
}
