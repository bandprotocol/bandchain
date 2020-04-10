#[macro_export]
macro_rules! prepare_entry_point {
    ($name:ident) => {
        #[no_mangle]
        pub fn prepare() {
            $name(BorshDeserialize::try_from_slice(&oei::get_calldata()).unwrap());
        }
    };
}

#[macro_export]
macro_rules! execute_entry_point {
    ($name:ident) => {
        #[no_mangle]
        pub fn execute() {
            oei::save_return_data(
                &$name(BorshDeserialize::try_from_slice(&oei::get_calldata()).unwrap())
                    .try_to_vec()
                    .unwrap(),
            );
        }
    };
}
