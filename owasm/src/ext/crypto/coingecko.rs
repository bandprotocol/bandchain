//! [CoinGecko.com](https://coingecko.com) Oracle Extension

use crate::core::{Oracle, ShellCmd};
use crate::ext::crypto::coins::Coins;
use crate::ext::utils::curl::Curl;

pub struct Price {
    symbol: String,
}

pub struct Volume24h {
    symbol: String,
}

impl Volume24h {
    pub fn new(coin: &Coins) -> Volume24h {
        Volume24h {
            symbol: String::from(match coin {
                Coins::ADA => "cardano",
                Coins::BAND => "band-protocol",
                Coins::BCH => "bitcoin-cash",
                Coins::BNB => "binance-coin",
                Coins::BTC => "bitcoin",
                Coins::EOS => "eos",
                Coins::ETH => "ethereum",
                Coins::LTC => "litecoin",
                Coins::ETC => "ethereum-classic",
                Coins::TRX => "tron",
                Coins::XRP => "ripple",
            }),
        }
    }
}

impl Price {
    pub fn new(coin: &Coins) -> Price {
        Price {
            symbol: String::from(match coin {
                Coins::ADA => "cardano",
                Coins::BAND => "band-protocol",
                Coins::BCH => "bitcoin-cash",
                Coins::BNB => "binance-coin",
                Coins::BTC => "bitcoin",
                Coins::EOS => "eos",
                Coins::ETH => "ethereum",
                Coins::LTC => "litecoin",
                Coins::ETC => "ethereum-classic",
                Coins::TRX => "tron",
                Coins::XRP => "ripple",
            }),
        }
    }
}

impl Oracle for Price {
    type T = f32;

    fn as_cmd(&self) -> ShellCmd {
        Curl::new(&[format!(
            "https://api.coingecko.com/api/v3/simple/price?ids={}&vs_currencies=usd",
            &self.symbol
        )])
        .as_cmd()
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        let parsed = json::parse(&output).ok()?;
        parsed[&self.symbol]["usd"].as_f32()
    }
}

impl Oracle for Volume24h {
    type T = f64;

    fn as_cmd(&self) -> ShellCmd {
        Curl::new(&[format!(
            "https://api.coingecko.com/api/v3/coins/{}/market_chart?vs_currency=usd&days=1",
            &self.symbol
        )])
        .as_cmd()
    }

    fn from_cmd_output(&self, output: String) -> Option<f64> {
        let parsed = json::parse(&output).ok()?;
        (parsed["total_volumes"].members().last()?)[1].as_f64()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new(&Coins::BTC).as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"]
            )
        );

        assert_eq!(
            Volume24h::new(&Coins::BTC).as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://api.coingecko.com/api/v3/coins/bitcoin/market_chart?vs_currency=usd&days=1"]
            )
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            Price::new(&Coins::BTC).from_cmd_output(r#"{"bitcoin":{"usd":100.0}}"#.into()),
            Some(100.0)
        );

        assert_eq!(
            Volume24h::new(&Coins::ETH).from_cmd_output(r#"{"prices":[[1579761912819,163.45762493199558],[1579762082019,163.49390724809746],[1579762372880,163.20299006856183]],"market_caps":[[1579761912819,17841379271.19409],[1579762082019,17841379271.19409],[1579762372880,17885078490.944096]],"total_volumes":[[1579761912819,8039481234.997162],[1579762082019,8812758589.411896],[1579762372880,8800955317.674072]]}"#.into()),
            Some(8800955317.674072)
        )
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Price::new(&Coins::BTC).from_cmd_output(r#"{}"#.into()), None);

        assert_eq!(
            Volume24h::new(&Coins::ETH).from_cmd_output(r#"{"prices":[[1579762372880,163.20299006856183]],"market_caps":[[1579762372880,17885078490.944096]],"total_volumes":[]}"#.into()),
            None
        )
    }

    #[test]
    fn test_request_all_tokens_price_from_coingecko() {
        println!(
            "{:?}",
            (vec![
                Price::new(&Coins::ADA),
                Price::new(&Coins::BAND),
                Price::new(&Coins::BCH),
                Price::new(&Coins::BNB),
                Price::new(&Coins::BTC),
                Price::new(&Coins::EOS),
                Price::new(&Coins::ETC),
                Price::new(&Coins::ETH),
                Price::new(&Coins::LTC),
                Price::new(&Coins::TRX),
                Price::new(&Coins::XRP),
            ])
            .iter()
            .map(|x| (x.symbol.as_str(), x.from_cmd_output(x.as_cmd().execute()).unwrap_or(0.0)))
            .collect::<Vec<(&str, f32)>>()
        );
    }

    #[test]
    fn test_request_all_tokens_volume24h_from_coingecko() {
        println!(
            "{:?}",
            (vec![
                Volume24h::new(&Coins::ADA),
                Volume24h::new(&Coins::BAND),
                Volume24h::new(&Coins::BCH),
                Volume24h::new(&Coins::BNB),
                Volume24h::new(&Coins::BTC),
                Volume24h::new(&Coins::EOS),
                Volume24h::new(&Coins::ETC),
                Volume24h::new(&Coins::ETH),
                Volume24h::new(&Coins::LTC),
                Volume24h::new(&Coins::TRX),
                Volume24h::new(&Coins::XRP),
            ])
            .iter()
            .map(|x| (x.symbol.as_str(), x.from_cmd_output(x.as_cmd().execute()).unwrap_or(0.0)))
            .collect::<Vec<(&str, f64)>>()
        );
    }
}
