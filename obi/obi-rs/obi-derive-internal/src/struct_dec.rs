use quote::quote;
use syn::export::TokenStream2;
use syn::{Fields, ItemStruct};

pub fn struct_dec(input: &ItemStruct) -> syn::Result<TokenStream2> {
    let name = &input.ident;
    let generics = &input.generics;
    let mut decode_field_types = TokenStream2::new();
    let return_value = match &input.fields {
        Fields::Named(fields) => {
            let mut body = TokenStream2::new();
            for field in &fields.named {
                let field_name = field.ident.as_ref().unwrap();
                let delta = {
                    let field_type = &field.ty;
                    decode_field_types.extend(quote! {
                        #field_type: obi::OBIDecode,
                    });
                    quote! {
                        #field_name: obi::OBIDecode::decode(buf)?,
                    }
                };
                body.extend(delta);
            }
            quote! {
                Self { #body }
            }
        }
        Fields::Unnamed(fields) => {
            let mut body = TokenStream2::new();
            for _ in 0..fields.unnamed.len() {
                let delta = quote! {
                    obi::OBIDecode::decode(buf)?,
                };
                body.extend(delta);
            }
            quote! {
                Self( #body )
            }
        }
        Fields::Unit => {
            quote! {
                Self {}
            }
        }
    };
    Ok(quote! {
        impl #generics obi::dec::OBIDecode for #name #generics where #decode_field_types {
            fn decode(buf: &mut &[u8]) -> std::result::Result<Self, std::io::Error> {
                Ok(#return_value)
            }
        }
    })
}
