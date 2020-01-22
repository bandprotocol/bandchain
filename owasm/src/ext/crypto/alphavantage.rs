//! [Alphavantage.co](https://www.alphavantage.co/) Oracle Extension

use crate::core::{Oracle, ShellCmd};
use crate::ext::crypto::coins::Coins;

pub struct Price {
    symbol: String,
    alphavantage_api_key: String,
}

impl Price {
    pub fn new(coin: &Coins, alphavantage_api_key: &String) -> Price {
        Price {
            symbol: String::from(match coin {
                Coins::BTC => "BTC",
                Coins::ETH => "ETH",
                Coins::BAND => "BAND",
            }),
            alphavantage_api_key: (*alphavantage_api_key).clone(),
        }
    }
}

impl Oracle for Price {
    type T = f32;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new(
            "curl",
            &[format!(
                "https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency={}&to_currency=USD&apikey={}",
                &self.symbol,
                &self.alphavantage_api_key
            )],
        )
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        let parsed = json::parse(&output).ok()?;
        Some(
            parsed["Realtime Currency Exchange Rate"]["5. Exchange Rate"]
                .as_str()?
                .parse::<f32>()
                .ok()?,
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new(&Coins::BTC, &"WVKPOO76169EX950".to_string()).as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://www.alphavantage.co/query?function=CURRENCY_EXCHANGE_RATE&from_currency=BTC&to_currency=USD&apikey=WVKPOO76169EX950"]
            )
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            Price::new(&Coins::BTC, &"".to_string()).from_cmd_output(r#"{"Realtime Currency Exchange Rate": {"1. From_Currency Code": "BAND","2. From_Currency Name": null,"3. To_Currency Code": "USD","4. To_Currency Name": "United States Dollar","5. Exchange Rate": "100.0","6. Last Refreshed": "2020-01-22 06:07:35","7. Time Zone": "UTC","8. Bid Price": "0.23960000","9. Ask Price": "0.24010000"}}"#.into()),
            Some(100.0)
        );
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Price::new(&Coins::BTC, &"".to_string()).from_cmd_output(r#"{}"#.into()), None);
        assert_eq!(
            Price::new(&Coins::BTC, &"just normal string".to_string())
                .from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new(&Coins::BTC, &"{}".to_string()).from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new(&Coins::BTC, &r#"{"UnknowField":"100.0"}"#.to_string())
                .from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new(
                &Coins::BTC,
                &r#"{"Realtime Currency Exchange Rate":"just normal string""#.to_string()
            )
            .from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new(&Coins::BTC, &r#"{"Realtime Currency Exchange Rate":{}}"#.to_string())
                .from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new(
                &Coins::BTC,
                &r#"{"Realtime Currency Exchange Rate":{"UnknowField":"100.0"}}"#.to_string()
            )
            .from_cmd_output(r#"{}"#.into()),
            None
        );
    }
}
