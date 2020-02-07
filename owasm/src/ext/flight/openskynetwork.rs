//! [opensky-network.org/](https://opensky-network.org/) Oracle Extension

use crate::core::{Oracle, ShellCmd};
use crate::ext::flight::flight_option::FlightOption;
use crate::ext::utils::curl::Curl;

pub struct Flight {
    flight_op: String,
    airport: String,
    icao24: String,
    should_happen_before: u64,
    should_happen_after: u64,
}

impl Flight {
    pub fn new(
        icao24: impl Into<String>,
        flight_op: &FlightOption,
        airport: impl Into<String>,
        should_happen_before: u64,
        should_happen_after: u64,
    ) -> Flight {
        Flight {
            flight_op: FlightOption::to_string(flight_op),
            airport: airport.into(),
            icao24: icao24.into(),
            should_happen_before,
            should_happen_after,
        }
    }
}

impl Oracle for Flight {
    type T = bool;

    fn as_cmd(&self) -> ShellCmd {
        Curl::new(&[format!(
            "https://opensky-network.org/api/flights/{}?airport={}&begin={}&end={}",
            &self.flight_op, &self.airport, &self.should_happen_after, &self.should_happen_before
        )])
        .as_cmd()
    }

    fn from_cmd_output(&self, output: String) -> Option<bool> {
        let parsed = json::parse(&output).ok()?;
        let mut members = parsed.members();
        if members.len() == 0 {
            return Some(false);
        }
        members.nth(0)?["icao24"].as_str()?;
        Some(
            members
                .filter(|x| x["icao24"].as_str() == Some(&self.icao24))
                .collect::<Vec<&json::JsonValue>>()
                .len()
                > 0,
        )
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            Flight::new("a79de9",&FlightOption::DEPARTURE, "EDDF", 1000, 900).as_cmd(),
            ShellCmd::new("curl", &["https://opensky-network.org/api/flights/departure?airport=EDDF&begin=900&end=1000"])
        );
    }

    #[test]
    fn test_from_cmd_found_flight_ok() {
        assert_eq!(
            Flight::new("a79de9",&FlightOption::DEPARTURE, "EDDF", 1580937440,1580997440).from_cmd_output(r#"[{"icao24":"406459","firstSeen":1580938171,"estDepartureAirport":"EDDF","lastSeen":1580941248,"estArrivalAirport":"EDDP","callsign":"BCS30C ","estDepartureAirportHorizDistance":2588,"estDepartureAirportVertDistance":194,"estArrivalAirportHorizDistance":7435,"estArrivalAirportVertDistance":132,"departureAirportCandidatesCount":1,"arrivalAirportCandidatesCount":6},{"icao24":"a79de9","firstSeen":1580937606,"estDepartureAirport":"EDDF","lastSeen":1580941006,"estArrivalAirport":"LFPG","callsign":"FDX36","estDepartureAirportHorizDistance":2854,"estDepartureAirportVertDistance":187,"estArrivalAirportHorizDistance":2411,"estArrivalAirportVertDistance":157,"departureAirportCandidatesCount":2,"arrivalAirportCandidatesCount":8},{"icao24":"3c5426","firstSeen":1580937574,"estDepartureAirport":"EDDF","lastSeen":1580940348,"estArrivalAirport":"EDDP","callsign":"BCS27H ","estDepartureAirportHorizDistance":2669,"estDepartureAirportVertDistance":194,"estArrivalAirportHorizDistance":4686,"estArrivalAirportVertDistance":25,"departureAirportCandidatesCount":2,"arrivalAirportCandidatesCount":6}]"#.into()),
            Some(true)
        );
    }

    #[test]
    fn test_from_cmd_not_found_flight_ok() {
        assert_eq!(
            Flight::new("a79de9",&FlightOption::DEPARTURE, "EDDF", 1580937440,1580997440).from_cmd_output(r#"[{"icao24":"406459","firstSeen":1580938171,"estDepartureAirport":"EDDF","lastSeen":1580941248,"estArrivalAirport":"EDDP","callsign":"BCS30C ","estDepartureAirportHorizDistance":2588,"estDepartureAirportVertDistance":194,"estArrivalAirportHorizDistance":7435,"estArrivalAirportVertDistance":132,"departureAirportCandidatesCount":1,"arrivalAirportCandidatesCount":6},{"icao24":"3c5426","firstSeen":1580937574,"estDepartureAirport":"EDDF","lastSeen":1580940348,"estArrivalAirport":"EDDP","callsign":"BCS27H ","estDepartureAirportHorizDistance":2669,"estDepartureAirportVertDistance":194,"estArrivalAirportHorizDistance":4686,"estArrivalAirportVertDistance":25,"departureAirportCandidatesCount":2,"arrivalAirportCandidatesCount":6}]"#.into()),
            Some(false)
        );
    }

    #[test]
    fn test_from_cmd_not_found_flight_for_empty_array_ok() {
        assert_eq!(
            Flight::new("a79de9", &FlightOption::DEPARTURE, "EDDF", 1580937440, 1580997440)
                .from_cmd_output(r#"[]"#.into()),
            Some(false)
        );
    }

    #[test]
    fn test_from_cmd_array_not_contain_object_not_ok() {
        assert_eq!(
            Flight::new("a79de9", &FlightOption::DEPARTURE, "EDDF", 1580997441, 1580997440)
                .from_cmd_output(
                    r#"["Start after end time or more than seven days of data requested"]"#.into()
                ),
            None
        );
    }

    #[test]
    fn test_from_cmd_not_array_not_ok() {
        assert_eq!(
            Flight::new("a79de9", &FlightOption::DEPARTURE, "EDDF", 1580997441, 1580997440)
                .from_cmd_output(
                    r#"The request arguments you provided were invalid. Please consult the documentation. Details: Failed to convert value of type 'java.lang.String' to required type 'int'; nested exception is java.lang.NumberFormatException: For input string: "999999999999999""#.into()
                ),
            None
        );
    }
}
