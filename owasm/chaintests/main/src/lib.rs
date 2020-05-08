use owasm::oei;

#[no_mangle]
pub fn prepare() {}

#[no_mangle]
pub fn execute() {
    let req_id = oei::get_request_id();
    let validators_count = oei::get_ask_count();

    oei::save_return_data(&vec![(req_id + validators_count) as u8, 0, 0, 0, 0, 0, 0, 0]);
}
