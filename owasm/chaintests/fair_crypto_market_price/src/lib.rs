use obi::{OBIDecode, OBIEncode};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode)]
struct Input {
    base_symbol: String,
    quote_symbol: String,
    aggregation_method: String,
    multiplier: u64,
}

#[derive(OBIEncode, Debug, PartialEq)]
struct Output {
    px: u64,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Coingecko data source
    oei::ask_external_data(10, 1, &input.base_symbol.as_bytes());
    oei::ask_external_data(11, 1, &input.quote_symbol.as_bytes());

    // Crypto compare source
    oei::ask_external_data(20, 2, &input.base_symbol.as_bytes());
    oei::ask_external_data(21, 2, &input.quote_symbol.as_bytes());

    // Binance source
    oei::ask_external_data(30, 3, &input.base_symbol.as_bytes());
    oei::ask_external_data(31, 3, &input.quote_symbol.as_bytes());
}

fn len(arr: &Vec<f64>) -> f64 {
    let mut l = 0f64;
    for _ in arr.iter() {
        l += 1f64
    }
    l
}

fn only_positive(arr: Vec<f64>) -> Vec<f64> {
    arr.iter().filter(|&&x| x > 0f64).map(|&x| x).collect::<Vec<_>>()
}

fn mean(arr: Vec<f64>) -> f64 {
    let pos_arr = only_positive(arr);
    let len_pos_arr = len(&pos_arr);
    if len_pos_arr > 0f64 {
        pos_arr.iter().fold(0f64, |sum, val| sum + val) / len_pos_arr
    } else {
        0f64
    }
}

fn median(arr: Vec<f64>) -> f64 {
    let mut pos_arr = only_positive(arr);
    let len_pos_arr = len(&pos_arr);
    if len_pos_arr > 0f64 {
        pos_arr.sort_by(|a, b| a.partial_cmp(b).unwrap());
        let mid = len_pos_arr / 2f64;
        if len_pos_arr as u64 % 2 == 0 {
            (pos_arr[(mid - 1f64) as usize] + pos_arr[mid as usize]) / 2f64
        } else {
            pos_arr[mid as usize]
        }
    } else {
        0f64
    }
}

fn aggregate(
    input: Input,
    b1: Vec<f64>,
    b2: Vec<f64>,
    b3: Vec<f64>,
    q1: Vec<f64>,
    q2: Vec<f64>,
    q3: Vec<f64>,
) -> Result<Output, String> {
    match &input.aggregation_method[..] {
        "mean" => {
            let b_avg = mean(vec![mean(b1), mean(b2), mean(b3)]);
            let q_avg = mean(vec![mean(q1), mean(q2), mean(q3)]);

            if q_avg <= 0f64 {
                Err(String::from("average of quote currency is negative"))
            } else {
                Ok(Output { px: ((b_avg * (input.multiplier as f64)) / q_avg) as u64 })
            }
        }
        "median" => {
            let b_med = median(vec![median(b1), median(b2), median(b3)]);
            let q_med = median(vec![median(q1), median(q2), median(q3)]);

            if q_med <= 0f64 {
                Err(String::from("median of quote currency is negative"))
            } else {
                Ok(Output { px: ((b_med * (input.multiplier as f64)) / q_med) as u64 })
            }
        }
        _ => Err(String::from("unknown method")),
    }
}

#[no_mangle]
fn execute_impl(input: Input) -> Output {
    let b1 = ext::load_input::<f64>(10);
    let q1 = ext::load_input::<f64>(11);

    let b2 = ext::load_input::<f64>(20);
    let q2 = ext::load_input::<f64>(21);

    let b3 = ext::load_input::<f64>(30);
    let q3 = ext::load_input::<f64>(31);

    aggregate(input, b1, b2, b3, q1, q2, q3).unwrap()
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_only_positive() {
        let expected1: Vec<f64> = vec![];
        assert_eq!(expected1, only_positive(vec![]));

        let expected2: Vec<f64> = vec![];
        assert_eq!(expected2, only_positive(vec![0f64]));

        let expected3: Vec<f64> = vec![7.000005];
        assert_eq!(expected3, only_positive(vec![0f64, 0f64, 7.000005]));

        let expected4: Vec<f64> = vec![0.00001, 1234.5678];
        assert_eq!(expected4, only_positive(vec![0f64, 0.00001, 0f64, -1f64, 1234.5678]));
    }

    #[test]
    fn test_mean() {
        assert_eq!(0., mean(vec![]));
        assert_eq!(5., mean(vec![5.]));
        assert_eq!(7.5, mean(vec![5., 10.]));
        assert_eq!(5.5, mean(vec![1., 2., 3., 4., 5., 6., 7., 8., 9., 10.]));
        assert_eq!(
            5.5,
            mean(vec![0., 0., 0., 1., 2., 3., 4., 0., 0., 5., 6., 7., 8., 9., 10., 0., 0., 0.])
        );
    }

    #[test]
    fn test_median() {
        assert_eq!(0., median(vec![]));
        assert_eq!(0., median(vec![0.]));
        assert_eq!(5., median(vec![5.]));
        assert_eq!(12.5, median(vec![15., 10.]));
        assert_eq!(5.5, median(vec![5., 7., 6., 9., 8., 10., 3., 1., 4., 2.]));
        assert_eq!(
            5.5,
            median(vec![0., 6., 0., 3., 0., 1., 0., 7., 10., 0., 0., 5., 9., 8., 2., 0., 0., 4.])
        );
    }

    fn create_mock_input(method: &str) -> Input {
        Input {
            base_symbol: String::from("BTC"),
            quote_symbol: String::from("USDT"),
            aggregation_method: String::from(method),
            multiplier: 10000,
        }
    }

    #[test]
    fn test_aggregate_unknown() {
        let unknown_input = create_mock_input("some_unknown_method");

        let x: Vec<f64> = vec![];
        assert_eq!(
            aggregate(
                unknown_input,
                x.clone(),
                x.clone(),
                x.clone(),
                x.clone(),
                x.clone(),
                x.clone()
            )
            .err(),
            Some(String::from("unknown method")),
        );
    }

    #[test]
    fn test_aggregate_mean() {
        let x: Vec<Vec<f64>> = vec![
            vec![100., 120.], // b1
            vec![90.],        // b2
            vec![],           // b3
            vec![1., 2., 3.], // q1
            vec![],           // q2
            vec![0., 3.],     // q3
        ];

        let mean_input = create_mock_input("mean");

        assert_eq!(
            aggregate(
                mean_input,
                x[0].clone(),
                x[1].clone(),
                x[2].clone(),
                x[3].clone(),
                x[4].clone(),
                x[5].clone()
            )
            .ok(),
            Some(Output { px: 400000 }),
        );
    }

    #[test]
    fn test_aggregate_median() {
        let x: Vec<Vec<f64>> = vec![
            vec![100., 130., 150., 200.], // b1
            vec![100.],                   // b2
            vec![],                       // b3
            vec![1., 2., 5., 10., 100.],  // q1
            vec![],                       // q2
            vec![0., 2.5],                // q3
        ];

        let median_input = create_mock_input("median");

        assert_eq!(
            aggregate(
                median_input,
                x[0].clone(),
                x[1].clone(),
                x[2].clone(),
                x[3].clone(),
                x[4].clone(),
                x[5].clone()
            )
            .ok(),
            Some(Output { px: 320000 }),
        );
    }
}
