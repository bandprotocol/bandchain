use crate::core::{Oracle, ShellCmd};

pub struct Curl {
    args: Vec<String>,
}

impl Curl {
    pub fn new(args: &[impl AsRef<str>]) -> Curl {
        Curl {
            args: args.into_iter().map(|x| x.as_ref().into()).collect(),
        }
    }
}

impl Oracle for Curl {
    type T = String;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new("curl", &self.args)
    }

    fn from_cmd_output(&self, output: String) -> Option<String> {
        Some(output)
    }
}
