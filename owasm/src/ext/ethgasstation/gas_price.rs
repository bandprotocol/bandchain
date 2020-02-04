use crate::core::{Oracle, ShellCmd};
use crate::ext::utils::curl::Curl;

pub struct Price {
    gas_option: String,
}

impl Price {
    pub fn new(gas_option: impl Into<String>) -> Price {
        Price { gas_option: gas_option.into() }
    }
}

impl Oracle for Price {
    type T = f32;

    fn as_cmd(&self) -> ShellCmd {
        Curl::new(&["https://ethgasstation.info/json/ethgasAPI.json"]).as_cmd()
    }

    fn from_cmd_output(&self, output: String) -> Option<f32> {
        Some((json::parse(&output).ok()?)[&self.gas_option].as_f32()? / 10.0)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Price::new("any").as_cmd(),
            ShellCmd::new("curl", &["https://ethgasstation.info/json/ethgasAPI.json"])
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            Price::new("fastest").from_cmd_output(r#"{"fast":80,"fastest":130,"safeLow":10,"average":20,"block_time":13.902757619738752,"blockNum":9385550,"speed":0.7777511075401192,"safeLowWait":9.8,"avgWait":2.1,"fastWait":0.5,"fastestWait":0.5,"gasPriceRange":{"4":231.7,"6":231.7,"8":231.7,"10":9.8,"15":5.9,"20":2.1,"25":1.8,"30":1.4,"35":1.4,"40":1.2,"45":1.2,"50":0.6,"55":0.6,"60":0.6,"65":0.6,"70":0.6,"75":0.6,"80":0.5,"85":0.5,"90":0.5,"95":0.5,"100":0.5,"105":0.5,"110":0.5,"115":0.5,"120":0.5,"125":0.5,"130":0.5}}"#.into()),
            Some(13.0)
        );
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(
            Price::new("fastest").from_cmd_output(r#"{"safeLow":10,"gasPriceRange":{}}"#.into()),
            None
        );
    }
}
