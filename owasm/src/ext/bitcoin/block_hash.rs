//! [blockcypher.com] (http://api.blockcypher.com/) Oracle Extension

use crate::core::{Oracle, ShellCmd};
use crate::ext::utils::curl::Curl;

pub struct Info {
    block_height: u64,
}

impl Info {
    pub fn new(block_height: u64) -> Info {
        Info { block_height }
    }
}

impl Oracle for Info {
    type T = [u64; 4];

    fn as_cmd(&self) -> ShellCmd {
        Curl::new(&[format!(
            "http://api.blockcypher.com/v1/btc/main/blocks/{}?txstart=1&limit=1",
            &self.block_height
        )])
        .as_cmd()
    }

    fn from_cmd_output(&self, output: String) -> Option<[u64; 4]> {
        let x = (json::parse(&output).ok()?)["hash"]
            .as_str()?
            .chars()
            .collect::<Vec<char>>()
            .chunks(16)
            .map(|c| u64::from_str_radix(&c.iter().collect::<String>(), 16).ok())
            .collect::<Option<Vec<u64>>>()?;
        match x.len() {
            4 => Some([x[0], x[1], x[2], x[3]]),
            _ => None,
        }
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Info::new(9999).as_cmd(),
            ShellCmd::new(
                "curl",
                &["http://api.blockcypher.com/v1/btc/main/blocks/9999?txstart=1&limit=1"]
            )
        );
    }

    #[test]
    fn test_from_cmd_ok() {
        assert_eq!(
            Info::new(9999).from_cmd_output(r#"{"hash":"00000000fbc97cc6c599ce9c24dd4a2243e2bfd518eda56e1d5e47d29e29c3a7","height":9999,"chain":"BTC.main","total":0,"fees":0,"size":216,"ver":1,"time":"2009-04-06T03:11:31Z","received_time":"2009-04-06T03:11:31Z","coinbase_addr":"","relayed_by":"","bits":486604799,"nonce":3568610608,"n_tx":1,"prev_block":"000000003dd32df94cfafd16e0a8300ea14d67dcfee9e1282786c2617b8daa09","mrkl_root":"5012c1d2a46d5684aa0331f0d8a900767c86c0fd83bb632f357b1ea11fa69179","txids":[],"depth":606033,"prev_block_url":"https://api.blockcypher.com/v1/btc/main/blocks/000000003dd32df94cfafd16e0a8300ea14d67dcfee9e1282786c2617b8daa09","tx_url":"https://api.blockcypher.com/v1/btc/main/txs/"}"#.into()),
            Some([4224285894, 14238638866937236002, 4891683067244946798, 2116207844832953255])
        );
    }

    #[test]
    fn test_from_cmd_not_ok() {
        assert_eq!(
            Info::new(9999999).from_cmd_output(r#"{"error": "Block 9999999 not found."}"#.into()),
            None
        );
    }
}
