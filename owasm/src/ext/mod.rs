//! # Owasm Standard Library
//!
//! TODO
use crate::oei;

pub fn load_average<T>(external_id: i64) -> T
where
    T: std::str::FromStr + num::Num + num::NumCast + std::marker::Copy,
{
    let mut vals: Vec<T> = Vec::new();
    for idx in 0..oei::get_requested_validator_count() {
        let external_data = oei::get_external_data(external_id, idx);
        match external_data.and_then(|x| x.parse::<T>().ok()) {
            Some(v) => vals.push(v),
            None => (),
        }
    }
    if vals.len() == 0 {
        T::zero()
    } else {
        average(&vals)
    }
}

fn average<T>(vals: &Vec<T>) -> T
where
    T: std::str::FromStr + num::Num + num::NumCast + std::marker::Copy,
{
    let mut sum = T::zero();
    for x in vals {
        sum = sum + *x;
    }
    sum / num::cast(vals.len()).unwrap()
}

pub fn load_majority<T>(external_id: i64) -> T
where
    T: std::str::FromStr,
{
    // TODO
    oei::get_external_data(external_id, 0).unwrap().parse::<T>().ok().unwrap()
}

pub fn load_median<T>(external_id: i64) -> Option<T>
where
    T: std::str::FromStr + std::cmp::PartialOrd + std::marker::Copy + num::Num + num::NumCast,
{
    let mut vals: Vec<T> = Vec::new();
    for idx in 0..oei::get_requested_validator_count() {
        let external_data = oei::get_external_data(external_id, idx);
        match external_data.and_then(|x| x.parse::<T>().ok()) {
            Some(v) => vals.push(v),
            None => (),
        }
    }

    if vals.len() == 0 {
        None
    } else {
        Some(median(&mut vals))
    }
}

fn median<T>(vals: &mut Vec<T>) -> T
where
    T: std::str::FromStr + std::cmp::PartialOrd + std::marker::Copy + num::Num + num::NumCast,
{
    vals.sort_by(|a, b| a.partial_cmp(b).unwrap());
    let mid = vals.len() / 2;
    if mid % 2 == 0 {
        (vals[mid - 1] + vals[mid]) / num::cast(2).unwrap()
    } else {
        vals[mid]
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_average_int() {
        let vals = vec![3, 2, 5, 7, 2, 9, 1];
        assert_eq!(average(&vals), 4);
    }

    #[test]
    fn test_average_float() {
        let vals = vec![3.0, 2.0, 5.0, 7.0, 2.0, 9.0, 1.0];
        assert_eq!(average(&vals), 4.142857142857143);
    }
    #[test]
    fn test_median_odd_int() {
        let mut vals = vec![3, 2, 5, 7, 2, 9, 1];
        assert_eq!(median(&mut vals), 3);
    }
    #[test]
    fn test_median_even_int() {
        let mut vals = vec![3, 2, 5, 7, 2, 10, 32, 1];
        assert_eq!(median(&mut vals), 4);

        let mut vals = vec![13, 36, 33, 45];
        assert_eq!(median(&mut vals), 34);
    }

    #[test]
    fn test_median_odd_float() {
        let mut vals = vec![3.5, 2.7, 5.1, 7.4, 2.0, 9.1, 1.9];
        assert_eq!(median(&mut vals), 3.5);
    }
    #[test]
    fn test_median_even_float() {
        let mut vals = vec![3.4, 2.0, 5.7, 7.1, 2.2, 10.1, 32.0, 1.8];
        assert_eq!(median(&mut vals), 4.55);

        let mut vals = vec![13.0, 36.0, 45.0, 33.0];
        assert_eq!(median(&mut vals), 34.5);
    }
}
