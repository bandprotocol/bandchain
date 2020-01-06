use crate::core::{Oracle, ShellCmd};

pub static BITCOIN: &str = "BTC";
pub static ETHEREUM: &str = "ETH";

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
                "https://min-api.cryptocompare.com/data/price?fsym={}&tsyms=USD",
                &self.symbol
            )],
        )
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        let parsed = json::parse(&output).ok()?;
        parsed["USD"].as_f32()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new("BTC").as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD"]
            )
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(Price::new("BTC").from_cmd_output(r#"{"USD":100.0}"#.into()), Some(100.0));
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Price::new("BTC").from_cmd_output(r#"{}"#.into()), None);
    }
}
