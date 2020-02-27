use owasm::oei;

#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[no_mangle]
pub fn prepare() {}

#[no_mangle]
pub fn execute() {
    let req_id = oei::get_request_id();
    let validators_count = oei::get_requested_validator_count();

    oei::save_return_data(&vec![(req_id + validators_count) as u8, 0, 0, 0, 0, 0, 0, 0]);
}
