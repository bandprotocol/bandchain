use crate::core::{Oracle, ShellCmd};

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
        ShellCmd::new("curl", &[format!("https://finance.yahoo.com/quote/{}", &self.symbol)])
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        let parsed = json::parse(
            &output
                .split("root.App.main =")
                .nth(1)?
                .split("(this)")
                .nth(0)?
                .split(";\n}")
                .nth(0)?,
        )
        .ok()?;
        parsed["context"]["dispatcher"]["stores"]["QuoteSummaryStore"]["price"]
            ["regularMarketPrice"]["raw"]
            .as_f32()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new("FB").as_cmd(),
            ShellCmd::new("curl", &["https://finance.yahoo.com/quote/FB"])
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(Price::new("FB").from_cmd_output(r#"... root.App.main = {"context":{"dispatcher":{"stores":{"QuoteSummaryStore":{"price":{"regularMarketPrice":{"raw":100.0}}}}}}}(this)) ...; "#.into()), Some(100.0));
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Price::new("FB").from_cmd_output(r#"{}"#.into()), None);
    }
}
