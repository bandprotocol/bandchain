//! [blockchain.info] (https://blockchain.info/q/getblockcount) Oracle Extension

use crate::core::{Oracle, ShellCmd};
use crate::ext::utils::curl::Curl;

pub struct Info {}

impl Info {
    pub fn new() -> Info {
        Info {}
    }
}

impl Oracle for Info {
    type T = u64;

    fn as_cmd(&self) -> ShellCmd {
        Curl::new(&["https://blockchain.info/q/getblockcount"]).as_cmd()
    }

    fn from_cmd_output(&self, output: String) -> Option<u64> {
        json::parse(&output).ok()?.as_u64()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Info::new().as_cmd(),
            ShellCmd::new("curl", &["https://blockchain.info/q/getblockcount"])
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(Info::new().from_cmd_output("616033".into()), Some(616033));
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Info::new().from_cmd_output("Block not found".into()), None);
    }
}
