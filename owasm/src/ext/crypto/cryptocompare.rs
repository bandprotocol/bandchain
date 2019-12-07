use crate::core::{Oracle, ShellCmd};

pub struct Price {}

impl Price {
    pub fn new() -> Price {
        Price {}
    }
}

impl Oracle for Price {
    type T = f32;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new("curl", &["https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD"])
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
            Price::new().as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD"]
            )
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(Price::new().from_cmd_output(r#"{"USD":100.0}"#.into()), Some(100.0));
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Price::new().from_cmd_output(r#"{}"#.into()), None);
    }
}
