use owasm::oei;

#[global_allocator]
static ALLOC: wee_alloc::WeeAlloc = wee_alloc::WeeAlloc::INIT;

#[no_mangle]
pub fn prepare() {
    oei::request_external_data(1, 1, "band-protocol");
}

#[no_mangle]
pub fn execute() {
    let data = oei::get_external_data(1, 0).unwrap();
    oei::save_return_data(&data.as_bytes());
}
