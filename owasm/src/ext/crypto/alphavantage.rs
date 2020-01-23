//! [Alphavantage.co](https://www.alphavantage.co/) Oracle Extension

use crate::core::{Oracle, ShellCmd};

pub struct Price {
    symbol: String,
    alphavantage_api_key: String,
}

impl Price {
    pub fn new(symbol: impl Into<String>, alphavantage_api_key: impl Into<String>) -> Price {
        Price { symbol: symbol.into(), alphavantage_api_key: alphavantage_api_key.into() }
    }
}

impl Oracle for Price {
    type T = f32;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new(
            "curl",
            &[format!(
                "https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol={}&apikey={}",
                &self.symbol, &self.alphavantage_api_key
            )],
        )
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        let parsed = json::parse(&output).ok()?;
        Some(parsed["Global Quote"]["05. price"].as_str()?.parse::<f32>().ok()?)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new("GOOG", "WVKPOO76169EX950").as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=GOOG&apikey=WVKPOO76169EX950"]
            )
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            Price::new("GOOG", "").from_cmd_output(r#"{"Global Quote":{"01. symbol":"GOOG","02. open":"1491.0000","03. high":"1503.2100","04. low":"1484.9300","05. price":"100.0","06. volume":"1593218","07. latest trading day":"2020-01-22","08. previous close":"1484.4000","09. change":"1.5500","10. change percent":"0.1044%"}}"#.into()),
            Some(100.0)
        );
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Price::new("GOOG", "").from_cmd_output(r#"{}"#.into()), None);
        assert_eq!(
            Price::new("GOOG", "just a normal string").from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(Price::new("GOOG", "{}").from_cmd_output(r#"{}"#.into()), None);
        assert_eq!(
            Price::new("GOOG", r#"{"UnknowField":"100.0"}"#).from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new("GOOG", r#"{"Realtime Currency Exchange Rate":"just normal string""#)
                .from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new("GOOG", r#"{"Realtime Currency Exchange Rate":{}}"#)
                .from_cmd_output(r#"{}"#.into()),
            None
        );
        assert_eq!(
            Price::new("GOOG", r#"{"Realtime Currency Exchange Rate":{"UnknowField":"100.0"}}"#)
                .from_cmd_output(r#"{}"#.into()),
            None
        );
    }
}
