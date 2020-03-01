use borsh::{BorshDeserialize, BorshSerialize};
use owasm::oei;
use owasm::{execute_entry_point, prepare_entry_point};

#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

const PRICE_DATA_SOURCE_ID: i64 = 1;
const PRICE_EXTERNAL_ID: i64 = 1;

#[derive(BorshDeserialize)]
struct Input {
    symbol: String,
    multiplier: u64,
}

#[derive(BorshSerialize)]
struct Output {
    px: u64,
}

fn load_average<T>(external_id: i64) -> T
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

fn load_majority<T>(external_id: i64) -> T
where
    T: std::str::FromStr,
{
    // TODO
    oei::get_external_data(external_id, 0).unwrap().parse::<T>().ok().unwrap()
}

fn prepare_impl(input: Input) {
    oei::request_external_data(PRICE_DATA_SOURCE_ID, PRICE_EXTERNAL_ID, input.symbol.as_bytes());
}

fn execute_impl(input: Input) -> Output {
    Output { px: (load_average::<f64>(PRICE_EXTERNAL_ID) * input.multiplier as f64) as u64 }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);

#[cfg(test)]
mod tests {

    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 5);
        // assert_eq!(
        //     shellwords::join(&["band protocol", "-f", "\"ezlife"]),
        //     "hello"
        // );
    }
}
