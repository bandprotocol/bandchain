use owasm::ext::utils::date;
use owasm::ext::weather::openweathermap;
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub city: String,
        pub key: String,
        pub sub_key: String,
    }
}

decl_data! {
    pub struct Data {
        pub weather_value: f64 = |params: &Parameter| openweathermap::WeatherInfo::new(&params.city, &params.key, &params.sub_key),
        pub time_stamp: u64 = |_: &Parameter| date::Date::new(),
    }
}

decl_result! {
    pub struct Result {
        pub weather_value: u64,
        pub time_stamp: u64,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    let mut total_weather_value = 0.0;
    let mut time_stamp_acc: u64 = 0;
    for each in &data {
        total_weather_value += each.weather_value;
        time_stamp_acc += each.time_stamp;
    }
    let average_weather_value = total_weather_value / (data.len() as f64);
    let avg_time_stamp = time_stamp_acc / (data.len() as u64);
    Result { weather_value: (average_weather_value * 100.0) as u64, time_stamp: avg_time_stamp }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter {
            city: String::from("seoul"),
            key: String::from("main"),
            sub_key: String::from("temp"),
        };
        let data1 = Data { weather_value: 100.0, time_stamp: 10 };
        let data2 = Data { weather_value: 200.0, time_stamp: 12 };
        assert_eq!(
            execute(params, vec![data1, data2]),
            Result { weather_value: 15000, time_stamp: 11 }
        );
    }

    #[test]
    fn test_call_real_price() {
        let params = Parameter {
            city: String::from("paris"),
            key: String::from("main"),
            sub_key: String::from("temp"),
        };
        let data = Data::build_from_local_env(&params).unwrap();
        println!(
            "Current temperature at Paris (times 100) is {:?} kelvin",
            execute(params, vec![data])
        );
    }
}
