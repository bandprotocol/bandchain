//! # Owasm Standard Library
//!
//! TODO

pub mod bitcoin;
pub mod crypto;
pub mod ethgasstation;
pub mod finance;
pub mod flight;
pub mod random;
pub mod utils;
pub mod weather;

use crate::oei;

pub fn load_average<T>(external_id: i64) -> T
where
    T: std::str::FromStr + num::Num,
{
    let mut count = T::zero();
    let mut sum = T::zero();
    for idx in 0..oei::get_requested_validator_count() {
        let external_data = oei::get_external_data(external_id, idx);
        match external_data.and_then(|x| x.parse::<T>().ok()) {
            Some(v) => {
                sum = sum + v;
                count = count + T::one();
            }
            None => (),
        }
    }
    sum / count
}

pub fn load_majority<T>(external_id: i64) -> T
where
    T: std::str::FromStr,
{
    // TODO
    oei::get_external_data(external_id, 0).unwrap().parse::<T>().ok().unwrap()
}
