use crate::core::{Oracle, ShellCmd};

pub struct Date {}

impl Date {
    pub fn new() -> Date {
        Date {}
    }
}

impl Oracle for Date {
    type T = u64;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new("date", &["+'%s'"])
    }

    fn from_cmd_output(&self, output: String) -> Option<u64> {
        output.parse::<u64>().ok()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(Date::new().as_cmd(), (ShellCmd::new("date", &["+'%s'"])));
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(Date::new().from_cmd_output("1579773294".into()), Some(1579773294));
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(Date::new().from_cmd_output("not integer".into()), None);
    }
}
