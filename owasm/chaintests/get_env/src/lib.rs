use owasm::oei;

#[no_mangle]
pub fn prepare() {}

#[no_mangle]
pub fn execute() {
    let calldata = oei::get_calldata();
    let fn_name = std::str::from_utf8(&calldata).unwrap();

    match fn_name {
        "getSufficientValidatorCount" => {
            let data = oei::get_sufficient_validator_count() as u64;
            oei::save_return_data(&data.to_be_bytes());
        }
        _ => oei::save_return_data(&[0u8; 8]),
    }
}
