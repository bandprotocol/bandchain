//! # Owasm Standard Library
//!
//! TODO
use crate::oei;

pub fn load_input<T>(external_id: i64) -> Vec<T>
where
    T: std::str::FromStr,
{
    let mut vals: Vec<T> = Vec::new();
    for idx in 0..oei::get_ask_count() {
        let external_data = oei::get_external_data(external_id, idx);
        match external_data.ok().and_then(|x| x.trim_end().parse::<T>().ok()) {
            Some(v) => vals.push(v),
            None => (),
        }
    }
    return vals;
}

pub fn load_average<T>(external_id: i64) -> T
where
    T: std::str::FromStr + num::Num,
{
    let vals = load_input(external_id);
    if vals.len() == 0 {
        T::zero()
    } else {
        average(vals)
    }
}

fn average<T>(vals: Vec<T>) -> T
where
    T: num::Num,
{
    let mut sum = T::zero();
    let mut count = T::zero();
    for x in vals {
        sum = sum + x;
        count = count + T::one();
    }
    sum / count
}

pub fn load_majority<T>(external_id: i64) -> Option<T>
where
    T: std::str::FromStr + std::cmp::PartialEq,
{
    let vec: Vec<T> = load_input(external_id);
    if vec.len() == 0 {
        None
    } else {
        majority(vec)
    }
}

fn majority<T>(mut vec: Vec<T>) -> Option<T>
where
    T: std::cmp::PartialEq,
{
    let mut candidate = 0;
    let mut count = 1;
    let len = vec.len();

    // Find majority by Boyer–Moore majority vote algorithm
    // https://en.wikipedia.org/wiki/Boyer%E2%80%93Moore_majority_vote_algorithm
    for idx in 1..len {
        if vec[candidate] == vec[idx] {
            count = count + 1;
        } else {
            count = count - 1;
        }
        if count == 0 {
            candidate = idx;
            count = 1;
        }
    }

    count = 0;
    for idx in 0..len {
        if vec[candidate] == vec[idx] {
            count = count + 1;
        }
    }

    if 2 * count > len {
        Some(vec.swap_remove(candidate))
    } else {
        None
    }
}

pub fn load_median<T>(external_id: i64) -> Option<T>
where
    T: std::str::FromStr + std::cmp::PartialOrd + num::Num,
{
    let vals = load_input(external_id);
    if vals.len() == 0 {
        None
    } else {
        Some(median(vals))
    }
}

fn median<T>(mut vals: Vec<T>) -> T
where
    T: std::cmp::PartialOrd + num::Num,
{
    vals.sort_by(|a, b| a.partial_cmp(b).unwrap());
    let mid = vals.len() / 2;
    if vals.len() % 2 == 0 {
        let rhs = vals.swap_remove(mid);
        let lhs = vals.swap_remove(mid - 1);
        (lhs + rhs) / (T::one() + T::one())
    } else {
        vals.swap_remove(mid)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_average_int() {
        let vals = vec![3, 2, 5, 7, 2, 9, 1];
        assert_eq!(average(vals), 4);
    }

    #[test]
    fn test_average_single_int() {
        let vals = vec![3];
        assert_eq!(average(vals), 3);
    }

    #[test]
    fn test_average_float() {
        let vals = vec![3.0, 2.0, 5.0, 7.0, 2.0, 9.0, 1.0];
        assert_eq!(average(vals), 4.142857142857143);
    }

    #[test]
    fn test_average_single_float() {
        let vals = vec![3.0];
        assert_eq!(average(vals), 3.0);
    }

    #[test]
    fn test_median_odd_int() {
        let vals = vec![3, 2, 5, 7, 2, 9, 1];
        assert_eq!(median(vals), 3);
    }

    #[test]
    fn test_median_single_int() {
        let vals = vec![3];
        assert_eq!(median(vals), 3);
    }

    #[test]
    fn test_median_even_int() {
        let vals = vec![3, 2, 5, 7, 2, 10, 32, 1];
        assert_eq!(median(vals), 4);
        let vals = vec![13, 36, 33, 45];
        assert_eq!(median(vals), 34);
        let vals = vec![13, 15];
        assert_eq!(median(vals), 14);
    }

    #[test]
    fn test_median_odd_float() {
        let vals = vec![3.5, 2.7, 5.1, 7.4, 2.0, 9.1, 1.9];
        assert_eq!(median(vals), 3.5);
    }

    #[test]
    fn test_median_single_float() {
        let vals = vec![3.0];
        assert_eq!(median(vals), 3.0);
    }

    #[test]
    fn test_median_even_float() {
        let vals = vec![3.4, 2.0, 5.7, 7.1, 2.2, 10.1, 32.0, 1.8];
        assert_eq!(median(vals), 4.55);
        let vals = vec![13.0, 36.0, 45.0, 33.0];
        assert_eq!(median(vals), 34.5);
        let vals = vec![13.0, 36.2];
        assert_eq!(median(vals), 24.6);
    }

    #[test]
    fn test_majority_int() {
        let vals = vec![1, 2, 3, 1, 3, 1, 1];
        assert_eq!(majority(vals), Some(1));
    }

    #[test]
    fn test_majority_single_int() {
        let vals = vec![3];
        assert_eq!(majority(vals), Some(3));
    }

    #[test]
    fn test_majority_int_result_none() {
        let vals = vec![1, 2, 3, 1, 3, 1, 1, 3];
        assert_eq!(majority(vals), None);
    }

    #[test]
    fn test_majority_float() {
        let vals = vec![0.3, 1.0, 0.3, 0.4, 1.0, 1.0, 1.0];
        assert_eq!(majority(vals), Some(1.0));
    }

    #[test]
    fn test_majority_single_float() {
        let vals = vec![3.0];
        assert_eq!(majority(vals), Some(3.0));
    }

    #[test]
    fn test_majority_float_result_none() {
        let vals = vec![0.3, 1.0, 0.3, 0.4, 4.0, 1.0, 99.99, 1.0, 1.0];
        assert_eq!(majority(vals), None);
    }

    #[test]
    fn test_majority_char() {
        let vals = vec!['a', 'b', 'a', 'b', 'b'];
        assert_eq!(majority(vals), Some('b'));
    }

    #[test]
    fn test_majority_single_char() {
        let vals = vec!['a'];
        assert_eq!(majority(vals), Some('a'));
    }

    #[test]
    fn test_majority_char_result_none() {
        let vals = vec!['a', 'b', 'a', 'b', 'c', 'b'];
        assert_eq!(majority(vals), None);
    }

    #[test]
    fn test_majority_string() {
        let vals = vec![String::from("mumu"), String::from("mumu"), String::from("momo")];
        assert_eq!(majority(vals), Some(String::from("mumu")));
    }

    #[test]
    fn test_majority_single_string() {
        let vals = vec![String::from("mumu")];
        assert_eq!(majority(vals), Some(String::from("mumu")));
    }

    #[test]
    fn test_majority_string_result_none() {
        let vals = vec![String::from("mumu"), String::from("momo")];
        assert_eq!(majority(vals), None);
    }
}
