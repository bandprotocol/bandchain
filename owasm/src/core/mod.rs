//! # Owasm Core Library
//!
//! TODO

mod cmd;
pub mod mem;
pub use cmd::ShellCmd;

pub trait Oracle {
    type T;

    fn as_cmd(&self) -> ShellCmd;
    fn from_cmd_output(&self, output: String) -> Option<Self::T>;
}

/// Encodes a list of shell commands to `Vec<u8>`.
pub fn encode_cmds(cmds: Vec<ShellCmd>) -> Option<Vec<u8>> {
    Some(serde_json::to_string(&cmds).ok()?.into_bytes())
}

/// Decodes an encoded list of shell commands back to the original form.
pub fn decode_cmds(raw: &[u8]) -> Option<Vec<ShellCmd>> {
    serde_json::from_slice(&raw).ok()
}

/// Encodes a list of shell command outputs to `Vec<u8>`.
pub fn encode_outputs(outputs: Vec<String>) -> Option<Vec<u8>> {
    Some(serde_json::to_string(&outputs).ok()?.into_bytes())
}

/// Decodes an encoded list of string outputs back to the original form.
pub fn decode_outputs(raw: &[u8]) -> Option<Vec<String>> {
    serde_json::from_slice(&raw).ok()
}

/// Executes the given list of commands and returns stdout outputs.
pub fn execute_with_local_env(cmds: Vec<ShellCmd>) -> Vec<String> {
    cmds.into_iter().map(|each| each.execute()).collect()
}
