//! [qrng.anu.edu.au](https://qrng.anu.edu.au/) Oracle Extension

use crate::core::{Oracle, ShellCmd};

pub struct RandomBytes {
    size: u8,
}

impl RandomBytes {
    pub fn new(size: u8) -> RandomBytes {
        RandomBytes { size }
    }
}

impl Oracle for RandomBytes {
    type T = Vec<u8>;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new(
            "curl",
            &[format!("https://qrng.anu.edu.au/API/jsonI.php?length={}&type=uint8", &self.size)],
        )
    }

    fn from_cmd_output(&self, output: String) -> Option<Vec<u8>> {
        let parsed = json::parse(&output).ok()?;
        let bytes = parsed["data"].members().map(|x| x.as_u8()).collect::<Option<Vec<u8>>>()?;
        if bytes.len() == self.size as usize {
            Some(bytes)
        } else {
            None
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            RandomBytes::new(10).as_cmd(),
            ShellCmd::new("curl", &["https://qrng.anu.edu.au/API/jsonI.php?length=10&type=uint8"])
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            RandomBytes::new(10).from_cmd_output(r#"{"type":"uint8","length":10,"data":[29,21,227,184,32,128,97,128,231,197],"success":true}"#.into()),
            Some(vec![29,21,227,184,32,128,97,128,231,197])
        );
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(RandomBytes::new(10).from_cmd_output(r#"{"success":false}"#.into()), None);
    }
}
