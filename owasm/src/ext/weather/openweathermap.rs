//! [openweathermap.org](https://api.openweathermap.org) Oracle Extension

use crate::core::{Oracle, ShellCmd};

pub struct WeatherInfo {
    city: String,
    key: String,
    sub_key: String,
}

impl WeatherInfo {
    pub fn new(
        city: impl Into<String>,
        key: impl Into<String>,
        sub_key: impl Into<String>,
    ) -> WeatherInfo {
        WeatherInfo { city: city.into(), key: key.into(), sub_key: sub_key.into() }
    }
}

impl Oracle for WeatherInfo {
    type T = f64;

    fn as_cmd(&self) -> ShellCmd {
        ShellCmd::new(
            "curl",
            &[format!(
                "https://api.openweathermap.org/data/2.5/weather?q={}&appid=ac7c05361f8f91652eab609377134ab7",
                &self.city
            )],
        )
    }

    fn from_cmd_output(&self, output: String) -> Option<f64> {
        let parsed = json::parse(&output).ok()?;
        if self.sub_key.is_empty() {
            return parsed[&self.key].as_f64();
        }
        parsed[&self.key][&self.sub_key].as_f64()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_as_cmd() {
        assert_eq!(
            WeatherInfo::new("seoul", "main", "temp").as_cmd(),
            ShellCmd::new(
                "curl",
                &["https://api.openweathermap.org/data/2.5/weather?q=seoul&appid=ac7c05361f8f91652eab609377134ab7"]
            )
        );
    }

    #[test]
    fn test_from_cmd_with_key_and_sub_key_ok() {
        assert_eq!(
            WeatherInfo::new("seoul", "main", "temp").from_cmd_output(r#"{"coord":{"lon":126.98,"lat":37.57},"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}],"base":"stations","main":{"temp":264.94,"feels_like":260.37,"temp_min":262.15,"temp_max":267.15,"pressure":1034,"humidity":44},"visibility":10000,"wind":{"speed":1.5,"deg":360},"clouds":{"all":1},"dt":1580905460,"sys":{"type":1,"id":8117,"country":"KR","sunrise":1580855589,"sunset":1580893130},"timezone":32400,"id":1835848,"name":"Seoul","cod":200}"#.into()),
            Some(264.94)
        );
    }

    #[test]
    fn test_from_cmd_with_key_only_ok() {
        assert_eq!(
            WeatherInfo::new("seoul", "visibility", "").from_cmd_output(r#"{"coord":{"lon":126.98,"lat":37.57},"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}],"base":"stations","main":{"temp":264.94,"feels_like":260.37,"temp_min":262.15,"temp_max":267.15,"pressure":1034,"humidity":44},"visibility":10000,"wind":{"speed":1.5,"deg":360},"clouds":{"all":1},"dt":1580905460,"sys":{"type":1,"id":8117,"country":"KR","sunrise":1580855589,"sunset":1580893130},"timezone":32400,"id":1835848,"name":"Seoul","cod":200}"#.into()),
            Some(10000.0)
        );
    }

    #[test]
    fn test_from_cmd_with_key_and_sub_key_not_ok() {
        assert_eq!(
            WeatherInfo::new("seoul", "main", "temp")
                .from_cmd_output(r#"{"cod":"404","message":"city not found"}"#.into()),
            None
        );

        assert_eq!(
            WeatherInfo::new("seoul", "sys", "country")
                .from_cmd_output(r#"{"coord":{"lon":126.98,"lat":37.57},"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}],"base":"stations","main":{"temp":264.94,"feels_like":260.37,"temp_min":262.15,"temp_max":267.15,"pressure":1034,"humidity":44},"visibility":10000,"wind":{"speed":1.5,"deg":360},"clouds":{"all":1},"dt":1580905460,"sys":{"type":1,"id":8117,"country":"KR","sunrise":1580855589,"sunset":1580893130},"timezone":32400,"id":1835848,"name":"Seoul","cod":200}"#.into()),
            None
        );
    }

    #[test]
    fn test_from_cmd_with_key_only_not_ok() {
        assert_eq!(
            WeatherInfo::new("seoul", "main", "").from_cmd_output(r#"{"coord":{"lon":126.98,"lat":37.57},"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}],"base":"stations","main":{"temp":264.94,"feels_like":260.37,"temp_min":262.15,"temp_max":267.15,"pressure":1034,"humidity":44},"visibility":10000,"wind":{"speed":1.5,"deg":360},"clouds":{"all":1},"dt":1580905460,"sys":{"type":1,"id":8117,"country":"KR","sunrise":1580855589,"sunset":1580893130},"timezone":32400,"id":1835848,"name":"Seoul","cod":200}"#.into()),
            None
        );

        assert_eq!(
            WeatherInfo::new("seoul", "name", "").from_cmd_output(r#"{"coord":{"lon":126.98,"lat":37.57},"weather":[{"id":800,"main":"Clear","description":"clear sky","icon":"01n"}],"base":"stations","main":{"temp":264.94,"feels_like":260.37,"temp_min":262.15,"temp_max":267.15,"pressure":1034,"humidity":44},"visibility":10000,"wind":{"speed":1.5,"deg":360},"clouds":{"all":1},"dt":1580905460,"sys":{"type":1,"id":8117,"country":"KR","sunrise":1580855589,"sunset":1580893130},"timezone":32400,"id":1835848,"name":"Seoul","cod":200}"#.into()),
            None
        );
    }
}
