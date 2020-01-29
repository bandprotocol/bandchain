use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, PartialEq)]
pub enum Coins {
    ADA,
    BAND,
    BCH,
    BNB,
    BTC,
    EOS,
    ETC,
    ETH,
    LTC,
    TRX,
    XRP,
}
