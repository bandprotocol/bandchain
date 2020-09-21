//! # Owasm Standard Library
use crate::oei;
pub mod cmp;
pub mod stats;

/// Returns an iterator of raw reports for the given external ID with nonzero status.
pub fn load_input_raw(eid: i64) -> impl Iterator<Item = String> {
    (0..oei::get_ask_count()).filter_map(move |idx| oei::get_external_data(eid, idx).ok())
}

/// Returns an iterator of raw data points for the given external ID, parsed into
/// the parameterized type using `std::str::FromStr` trait. Skip data points
/// with nonzero status OR cannot be parsed.
pub fn load_input<T>(eid: i64) -> impl Iterator<Item = T>
where
    T: std::str::FromStr,
{
    load_input_raw(eid).filter_map(|e| e.trim_end().parse::<T>().ok())
}

/// Returns the average value of the given external ID, ignoring unsuccessful reports.
pub fn load_average<T>(eid: i64) -> Option<T>
where
    T: std::str::FromStr + num::Num,
{
    stats::average(load_input(eid).collect())
}

/// Returns the median value of the given external ID, ignoring unsuccessful reports.
pub fn load_median<T>(eid: i64) -> Option<T>
where
    T: std::str::FromStr + std::cmp::Ord + num::Num,
{
    stats::median(load_input(eid).collect())
}

/// Returns the majority value of the given external ID, ignoring unsuccessful reports.
pub fn load_majority<T>(eid: i64) -> Option<T>
where
    T: std::str::FromStr + std::cmp::PartialEq,
{
    stats::majority(load_input(eid).collect())
}
