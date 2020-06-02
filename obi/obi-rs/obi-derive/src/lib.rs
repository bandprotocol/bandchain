extern crate proc_macro;
use obi_derive_internal::*;
use proc_macro::TokenStream;
use syn::ItemStruct;

#[proc_macro_derive(OBIEncode)]
pub fn obi_encode(input: TokenStream) -> TokenStream {
    let res = if let Ok(input) = syn::parse::<ItemStruct>(input.clone()) {
        struct_enc(&input)
    } else {
        // Derive macros can only be defined on structs
        unreachable!()
    };
    TokenStream::from(match res {
        Ok(res) => res,
        Err(err) => err.to_compile_error(),
    })
}

#[proc_macro_derive(OBIDecode)]
pub fn obi_decode(input: TokenStream) -> TokenStream {
    let res = if let Ok(input) = syn::parse::<ItemStruct>(input.clone()) {
        struct_dec(&input)
    } else {
        // Derive macros can only be defined on structs.
        unreachable!()
    };
    TokenStream::from(match res {
        Ok(res) => res,
        Err(err) => err.to_compile_error(),
    })
}
