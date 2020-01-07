/// # Macro to define data types that support Owasm encoding
///
/// By using this macro, the nested struct automatically implements three functions.
/// 1. `__input` static method that returns a list of commands required to build the data.
/// 2. `__output` static method that parse a list command outputs into the data.
/// 3. `build_from_local_env` static method to build data from local environment.
///
/// ## Examples
///
/// The following code snippet defines `Data` struct that consists of two `f32` fields: `coin_gecko` for
/// getting Bitcoin price from CoinGecko, and `crypto_compare` for getting Bitcoin price from CryptoCompare.
///
/// ```
/// use owasm::decl_data;
/// use owasm::ext::crypto::{coingecko, cryptocompare};
///
/// pub struct __Params {
///     pub symbol_cg: String,
///     pub symbol_cc: String,
/// }
/// decl_data! {
///     pub struct Data {
///         pub coin_gecko: f32 = |params: &__Params| coingecko::Price::new(&params.symbol_cg),
///         pub crypto_compare: f32 = |params: &__Params| cryptocompare::Price::new(&params.symbol_cc),
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
            pub fn __input(params: &__Params) -> Vec<ShellCmd> {
                vec![ $($field_howto(&params).as_cmd(),)* ]
            }

            pub fn __output(params: &__Params, mut output: Vec<String>) -> Option<$data_name> {
                Some($data_name {
                    $($field_name : $field_howto(&params).from_cmd_output(output.remove(0))?,)*
                })
            }

            pub fn build_from_local_env(params: &__Params) -> Option<$data_name> {
                Self::__output(params, execute_with_local_env(Self::__input(params)))
            }

            pub fn fields() -> Vec<String> {
                vec![ $(String::from(stringify!($field_name)),)*]
            }
        }

        pub type __Data = $data_name;
    };
}

#[macro_export]
macro_rules! decl_params {
    ($(pub $field_name:ident : $field_type:ty,)*) => {
        use serde::{Deserialize, Serialize};

        #[derive(Debug, Serialize, Deserialize)]
        pub struct __Params {
            $(pub $field_name : $field_type,)*
        }

        impl __Params {
            pub fn fields() -> Vec<String> {
                vec![ $(String::from(stringify!($field_name)),)*]
            }
        }
    };
}
