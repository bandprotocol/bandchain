use owasm::ext::flight::{flight_option, openskynetwork};
use owasm::{decl_data, decl_params, decl_result};

decl_params! {
    pub struct Parameter {
        pub icao24: String,
        pub flight_option: flight_option::FlightOption,
        pub airport: String,
        pub should_happen_before: u64,
        pub should_happen_after: u64,
    }
}

decl_data! {
    pub struct Data {
        pub has_flight_found: bool = |params: &Parameter| openskynetwork::Flight::new(
            &params.icao24,
            &params.flight_option,
            &params.airport,
            params.should_happen_before,
            params.should_happen_after,
        ),
    }
}

decl_result! {
    pub struct Result {
        pub has_flight_found: bool,
    }
}

pub fn execute(_params: Parameter, data: Vec<Data>) -> Result {
    let mut found_count = 0;
    let mut not_found_count = 0;
    for each in &data {
        if each.has_flight_found {
            found_count += 1;
        } else {
            not_found_count += 1;
        }
    }
    Result { has_flight_found: found_count > not_found_count }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_execute() {
        let params = Parameter {
            icao24: String::from("a79de9"),
            flight_option: flight_option::FlightOption::DEPARTURE,
            airport: String::from("EDDF"),
            should_happen_after: 1580937440,
            should_happen_before: 1580997440,
        };
        let data1 = Data { has_flight_found: true };
        let data2 = Data { has_flight_found: true };
        let data3 = Data { has_flight_found: false };
        assert_eq!(execute(params, vec![data1, data2, data3]), Result { has_flight_found: true });
    }

    #[test]
    fn test_call_and_finding_flight() {
        let params = Parameter {
            icao24: String::from("a79de9"),
            flight_option: flight_option::FlightOption::DEPARTURE,
            airport: String::from("EDDF"),
            should_happen_after: 1580937440,
            should_happen_before: 1580997440,
        };
        let data = Data::build_from_local_env(&params).unwrap();
        println!(
            "Is there a flight with icao24 that departure at EDDF after 1580937440 and before 1580997440 ? \n Answer is {:?}",
            execute(params, vec![data])
        );
    }
}
