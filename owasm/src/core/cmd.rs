use serde::{Deserialize, Serialize};
use std::process::Command;

#[derive(Debug, Serialize, Deserialize, PartialEq, Eq)]
pub struct ShellCmd {
    cmd: String,
    args: Vec<String>,
}

impl ShellCmd {
    pub fn new(cmd: impl AsRef<str>, args: &[impl AsRef<str>]) -> ShellCmd {
        ShellCmd {
            cmd: cmd.as_ref().into(),
            args: args.into_iter().map(|x| x.as_ref().into()).collect(),
        }
    }

    pub fn execute(&self) -> String {
        String::from_utf8(Command::new(&self.cmd).args(&self.args).output().unwrap().stdout)
            .unwrap()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cmd_input() {
        assert_eq!(
      serde_json::to_string(&ShellCmd::new(
        "curl",
        &["https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"],
      ))
      .unwrap(),
      r#"{"cmd":"curl","args":["https://api.coingecko.com/api/v3/simple/price?ids=bitcoin&vs_currencies=usd"]}"#
    );
    }
}
