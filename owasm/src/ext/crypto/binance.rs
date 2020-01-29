//! [Binance.com] (https://www.binance.com) Oracle Extension

use crate::core::{Oracle, ShellCmd};
use crate::ext::crypto::coins::Coins;

pub struct Price {
    symbol: String,
}

impl Price {
    pub fn new(coin: &Coins) -> Price {
        Price {
            symbol: String::from(match coin {
                Coins::ADA => "ADA",
                Coins::BAND => "BAND",
                Coins::BCH => "BCH",
                Coins::BNB => "BNB",
                Coins::BTC => "BTC",
                Coins::EOS => "EOS",
                Coins::ETH => "ETH",
                Coins::LTC => "LTC",
                Coins::ETC => "ETC",
                Coins::TRX => "TRX",
                Coins::XRP => "XRP",
            }),
        }
    }
}

impl Oracle for Price {
    type T = f32;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new(
            "curl",
            &[format!("https://api.binance.com/api/v1/depth?symbol={}USDT&limit=5", &self.symbol)],
        )
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        let parsed = json::parse(&output).ok()?;
        Some(
            (parsed["asks"][0][0].as_str()?.parse::<f32>().ok()?
                + parsed["bids"][0][0].as_str()?.parse::<f32>().ok()?)
                / 2.0,
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new(&Coins::BTC).as_cmd(),
            ShellCmd::new("curl", &["https://api.binance.com/api/v1/depth?symbol=BTCUSDT&limit=5"])
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            Price::new(&Coins::BTC).from_cmd_output(r#"{"lastUpdateId":19363463,"bids":[["100.0","12646.98000000"],["0.23500000","13998.54000000"],["0.23430000","821.59000000"],["0.23420000","4931.68000000"],["0.23410000","72.45000000"]],"asks":[["120.0","289.00000000"],["0.23750000","326.08000000"],["0.23760000","143.78000000"],["0.23770000","1003.61000000"],["0.23790000","49.91000000"]]}"#.into()),
            Some(110.0)
        );
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(
            Price::new(&Coins::BTC).from_cmd_output(
                r#"{"lastUpdateId":19363463,"bids":["100.0"],"asks":["120.0"]}"#.into()
            ),
            None
        );
    }

    #[test]
    fn test_request_all_tokens_price_from_binance() {
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
}
