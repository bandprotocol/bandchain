use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize, PartialEq)]
pub enum FlightOption {
    ARRIVAL,
    DEPARTURE,
}

impl FlightOption {
    pub fn to_string(flight_op: &FlightOption) -> String {
        match flight_op {
            FlightOption::ARRIVAL => String::from("arrival"),
            FlightOption::DEPARTURE => String::from("departure"),
        }
    }
}
