/// # Macro to define data types that support Owasm encoding
///
/// By using this macro, the nested struct automatically implements three functions.
/// 1. `__input` static method that returns a list of commands required to build the data.
/// 2. `__output` static method that parse a list command outputs into the data.
/// 3. `build_from_local_env` static method to build data from local environment.
///
/// ## Examples
///
/// The following code snippet defines `Data` struct that consists of two `f32` fields: `cgk` for
/// getting Bitcoin price from CoinGecko, and `ccc` for getting Bitcoin price from CryptoCompare.
///
/// ```
/// use owasm::decl_data;
/// use owasm::ext::crypto::{coingecko, cryptocompare};
///
/// decl_data! {
///     pub struct Data {
///         pub cgk: f32 = coingecko::Price::new("bitcoin"),
///         pub ccc: f32 = cryptocompare::Price::new(),
///     }
/// }
/// ```
#[macro_export]
macro_rules! decl_data {
    (pub struct $data_name:ident {
        $(pub $field_name:ident : $field_type:ty = $field_howto:expr ,)*
    }) => {
        use $crate::core::{Oracle, ShellCmd, execute_with_local_env};

        #[derive(Debug)]
        pub struct $data_name {
            $(pub $field_name : $field_type,)*
        }

        impl $data_name {
            pub fn __input() -> Vec<ShellCmd> {
                vec![ $($field_howto.as_cmd(),)* ]
            }

            pub fn __output(mut output: Vec<String>) -> Option<$data_name> {
                Some($data_name {
                    $($field_name : $field_howto.from_cmd_output(output.remove(0))?,)*
                })
            }

            pub fn build_from_local_env() -> Option< $data_name > {
                Self::__output(execute_with_local_env(Self::__input()))
            }
        }

        pub type __Data = $data_name;
    };
}
