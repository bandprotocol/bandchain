//! [CoinGecko.com](https://coingecko.com) Oracle Extension

use crate::core::{Oracle, ShellCmd};

pub static BITCOIN: &str = "bitcoin";
pub static ETHEREUM: &str = "ethereum";

pub struct Price {
    symbol: String,
}

impl Price {
    pub fn new(symbol: impl Into<String>) -> Price {
        Price { symbol: symbol.into() }
    }
}

impl Oracle for Price {
    type T = f32;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new(
            "curl",
            &[format!(
                "https://api.coingecko.com/api/v3/simple/price?ids={}&vs_currencies=usd",
                &self.symbol
            )],
        )
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        let parsed = json::parse(&output).ok()?;
        parsed[&self.symbol]["usd"].as_f32()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new("bitcoin").as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"]
            )
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            Price::new("bitcoin").from_cmd_output(r#"{"bitcoin":{"usd":100.0}}"#.into()),
            Some(100.0)
        );
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Price::new("bitcoin").from_cmd_output(r#"{}"#.into()), None);
    }
}
