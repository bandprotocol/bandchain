use core::cmp::Ordering;
use num::Float;

/// A comparison function on Float data types that work with NaN.
pub fn fcmp<T>(lhs: &T, rhs: &T) -> Ordering
where
    T: Float,
{
    match lhs.partial_cmp(&rhs) {
        Some(ordering) => ordering,
        None => {
            if lhs.is_nan() {
                if rhs.is_nan() {
                    Ordering::Equal
                } else {
                    Ordering::Greater
                }
            } else {
                Ordering::Less
            }
        }
    }
}
