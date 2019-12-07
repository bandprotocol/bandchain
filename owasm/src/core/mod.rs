//! # Owasm I/O Core Library
//!
//! TODO

mod cmd;

pub use cmd::ShellCmd;

pub trait Oracle {
    type T;

    fn as_cmd(&self) -> ShellCmd;
    fn from_cmd_output(&self, output: String) -> Option<Self::T>;
}

pub fn encode_cmds(cmds: Vec<ShellCmd>) -> Option<Vec<u8>> {
    Some(serde_json::to_string(&cmds).ok()?.into_bytes())
}

pub fn decode_cmds(raw: &[u8]) -> Option<Vec<ShellCmd>> {
    serde_json::from_slice(&raw).ok()
}

pub fn encode_outputs(outputs: Vec<String>) -> Option<Vec<u8>> {
    Some(serde_json::to_string(&outputs).ok()?.into_bytes())
}

pub fn decode_outputs(raw: &[u8]) -> Option<Vec<String>> {
    serde_json::from_slice(&raw).ok()
}

pub fn execute_with_local_env(cmds: Vec<ShellCmd>) -> Vec<String> {
    cmds.into_iter().map(|each| each.execute()).collect()
}
