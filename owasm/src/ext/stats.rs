use core::cmp::{Ord, Ordering, PartialEq};
use num::Num;

/// Returns the average value of the given data set, or None if data is empty.
pub fn average<T>(data: Vec<T>) -> Option<T>
where
    T: Num,
{
    let mut sum = T::zero();
    let mut count = T::zero();
    for v in data {
        sum = sum + v;
        count = count + T::one();
    }
    if count == T::zero() {
        None
    } else {
        Some(sum / count)
    }
}

/// Returns the median value using the given compare function, or None if data is empty.
pub fn median_by<T, F>(mut data: Vec<T>, compare: F) -> Option<T>
where
    T: Num,
    F: FnMut(&T, &T) -> Ordering,
{
    if data.len() == 0 {
        return None;
    }
    data.sort_by(compare);
    let mid = data.len() / 2;
    if data.len() % 2 == 0 {
        let rhs = data.swap_remove(mid);
        let lhs = data.swap_remove(mid - 1);
        Some((lhs + rhs) / (T::one() + T::one()))
    } else {
        Some(data.swap_remove(mid))
    }
}

/// Returns the median value of the given data set, or None if data is empty.
pub fn median<T>(data: Vec<T>) -> Option<T>
where
    T: Ord + Num,
{
    median_by(data, T::cmp)
}

/// Returns the majority value of the given data set, or None if there is no majority.
pub fn majority<T>(mut data: Vec<T>) -> Option<T>
where
    T: PartialEq,
{
    let mut candidate = 0;
    let mut count = 1;
    let len = data.len();

    // Find majority by Boyer–Moore majority vote algorithm
    // https://en.wikipedia.org/wiki/Boyer%E2%80%93Moore_majority_vote_algorithm
    for idx in 1..len {
        if data[candidate] == data[idx] {
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
        if data[candidate] == data[idx] {
            count = count + 1;
        }
    }

    if 2 * count > len {
        Some(data.swap_remove(candidate))
    } else {
        None
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::ext::cmp;

    #[test]
    fn test_average_empty() {
        let vals: Vec<i64> = vec![];
        assert_eq!(average(vals), None);
    }

    #[test]
    fn test_average_int() {
        let vals = vec![3, 2, 5, 7, 2, 9, 1];
        assert_eq!(average(vals), Some(4));
    }

    #[test]
    fn test_average_single_int() {
        let vals = vec![3];
        assert_eq!(average(vals), Some(3));
    }

    #[test]
    fn test_average_float() {
        let vals = vec![3.0, 2.0, 5.0, 7.0, 2.0, 9.0, 1.0];
        assert_eq!(average(vals), Some(4.142857142857143));
    }

    #[test]
    fn test_average_single_float() {
        let vals = vec![3.0];
        assert_eq!(average(vals), Some(3.0));
    }

    #[test]
    fn test_median_odd_int() {
        let vals = vec![3, 2, 5, 7, 2, 9, 1];
        assert_eq!(median(vals), Some(3));
    }

    #[test]
    fn test_median_single_int() {
        let vals = vec![3];
        assert_eq!(median(vals), Some(3));
    }

    #[test]
    fn test_median_empty() {
        let vals: Vec<i64> = vec![];
        assert_eq!(median(vals), None);
    }

    #[test]
    fn test_median_even_int() {
        let vals = vec![3, 2, 5, 7, 2, 10, 32, 1];
        assert_eq!(median(vals), Some(4));
        let vals = vec![13, 36, 33, 45];
        assert_eq!(median(vals), Some(34));
        let vals = vec![13, 15];
        assert_eq!(median(vals), Some(14));
    }

    #[test]
    fn test_median_odd_float() {
        let vals = vec![3.5, 2.7, 5.1, 7.4, 2.0, 9.1, 1.9];
        assert_eq!(median_by(vals, cmp::fcmp), Some(3.5));
    }

    #[test]
    fn test_median_single_float() {
        let vals = vec![3.0];
        assert_eq!(median_by(vals, cmp::fcmp), Some(3.0));
    }

    #[test]
    fn test_median_even_float() {
        let vals = vec![3.4, 2.0, 5.7, 7.1, 2.2, 10.1, 32.0, 1.8];
        assert_eq!(median_by(vals, cmp::fcmp), Some(4.55));
        let vals = vec![13.0, 36.0, 45.0, 33.0];
        assert_eq!(median_by(vals, cmp::fcmp), Some(34.5));
        let vals = vec![13.0, 36.2];
        assert_eq!(median_by(vals, cmp::fcmp), Some(24.6));
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
