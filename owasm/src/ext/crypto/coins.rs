use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, PartialEq)]
pub enum Coins {
    BTC,
    ETH,
    BAND,
}
