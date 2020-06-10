use hex;
use obi::{get_schema, OBIDecode, OBIEncode, OBISchema};
use owasm::{execute_entry_point, ext, oei, prepare_entry_point};

#[derive(OBIDecode, OBISchema)]
struct Input {
    size: u64,
}

#[derive(OBIEncode, OBISchema)]
struct Output {
    random_bytes: String,
}

#[no_mangle]
fn prepare_impl(input: Input) {
    // Quantum random data source
    oei::ask_external_data(1, 13, input.size.to_string().as_bytes());
}

fn accumulate_hex_strings(strings: Vec<String>, input_size: usize) -> String {
    hex::encode(
        strings
            .iter()
            .map(|x1| x1.split(",").map(|x2| x2.parse::<u8>().unwrap()).collect::<Vec<_>>())
            .fold(vec![0; input_size], |acc, x1| {
                acc.iter().zip(x1.iter()).map(|x2| x2.0 ^ x2.1).collect::<Vec<_>>()
            }),
    )
}

#[no_mangle]
fn execute_impl(input: Input) -> Output {
    Output {
        random_bytes: accumulate_hex_strings(ext::load_input::<String>(1), input.size as usize),
    }
}

prepare_entry_point!(prepare_impl);
execute_entry_point!(execute_impl);

#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::*;

    #[test]
    fn test_get_schema() {
        let mut schema = HashMap::new();
        Input::add_definitions_recursively(&mut schema);
        Output::add_definitions_recursively(&mut schema);
        let input_schema = get_schema(String::from("Input"), &schema);
        let output_schema = get_schema(String::from("Output"), &schema);
        println!("{}/{}", input_schema, output_schema);
        assert_eq!(
            "{size:u64}/{random_bytes:string}",
            format!("{}/{}", input_schema, output_schema),
        );
    }

    #[test]
    fn test_accumulate_hex_strings_1_1() {
        let hex_string = accumulate_hex_strings(vec![String::from("120")], 1);
        assert_eq!("78", hex_string)
    }

    #[test]
    fn test_accumulate_hex_strings_1_4() {
        let hex_string = accumulate_hex_strings(
            vec![String::from("1"), String::from("26"), String::from("100"), String::from("243")],
            1,
        );
        assert_eq!("8c", hex_string)
    }

    #[test]
    fn test_accumulate_hex_strings_1_10() {
        let hex_string = accumulate_hex_strings(
            vec![
                String::from("153"),
                String::from("106"),
                String::from("55"),
                String::from("209"),
                String::from("44"),
                String::from("32"),
                String::from("97"),
                String::from("102"),
                String::from("95"),
                String::from("224"),
            ],
            1,
        );
        assert_eq!("a1", hex_string)
    }

    #[test]
    fn test_accumulate_hex_strings_4_5() {
        let hex_string = accumulate_hex_strings(
            vec![
                String::from("1,2,3,4"),
                String::from("4,3,2,1"),
                String::from("100,43,151,255"),
                String::from("236,7,45,180"),
                String::from("213,55,87,56"),
            ],
            4,
        );
        assert_eq!("581aec76", hex_string);
    }

    #[test]
    fn test_accumulate_hex_strings_100_4() {
        let hex_string = accumulate_hex_strings(
            vec![
                String::from("154,79,163,42,188,109,142,139,176,65,65,169,204,168,226,224,127,217,162,224,169,215,18,17,192,143,180,245,196,93,103,240,126,247,44,194,177,135,160,100,131,148,120,26,183,147,8,102,104,191,115,135,76,11,67,212,161,18,117,83,234,36,83,0,54,65,194,99,200,84,21,214,158,42,75,100,225,117,117,95,88,61,85,200,154,202,219,209,95,190,132,22,112,165,65,4,92,189,141,170"),
                String::from("176,123,57,16,26,196,215,48,240,151,210,239,33,87,169,139,120,18,16,208,236,237,2,222,216,88,18,145,107,55,127,119,163,22,20,101,167,159,174,239,77,122,231,20,191,124,42,6,37,0,80,33,17,148,57,191,172,241,137,36,24,128,55,36,217,251,241,106,211,71,185,64,127,192,52,108,209,52,249,81,134,189,48,154,93,193,241,85,212,18,60,81,23,204,180,64,57,35,35,69"),
                String::from("127,70,0,4,187,86,137,29,206,79,194,119,58,47,201,140,4,61,229,199,83,72,48,145,237,113,31,79,202,172,141,66,34,175,17,161,45,13,123,8,82,221,31,116,220,60,234,225,197,0,202,78,25,120,159,117,177,39,199,88,234,151,251,224,28,150,199,213,146,214,25,113,178,157,199,188,212,182,151,51,29,255,109,121,200,137,15,161,33,140,249,15,78,104,110,167,250,68,228,63"),
                String::from("5,93,181,114,150,64,214,122,33,163,48,126,110,25,0,44,53,41,3,44,181,197,193,35,104,207,63,47,227,90,222,47,189,44,187,120,52,118,68,224,184,252,185,153,238,123,149,150,87,223,51,220,10,47,126,171,102,200,208,146,143,91,54,243,82,213,226,76,209,185,0,201,116,190,148,178,121,66,247,82,85,226,90,151,198,99,48,154,14,11,105,19,206,186,176,9,189,225,83,13"),
            ],
            100,
        );
        assert_eq!(
            "502f2f4c8bbf06dcaf3a614fb9c982cb36df54dba3b7e17d9d698604869c4bea4262927e0f63316324cf39e33aa85d17df60da344ec89bb5da0cebbd9768a937a1f91690587cb52e27c92c069db5ec6f969d52bcc9e115bfa42b285be7bb2bea223b19dd",
            hex_string
        );
    }
}
