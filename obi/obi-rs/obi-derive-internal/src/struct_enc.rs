use quote::quote;
use syn::export::{Span, TokenStream2};
use syn::{Fields, Index, ItemStruct};

pub fn struct_enc(input: &ItemStruct) -> syn::Result<TokenStream2> {
    let name = &input.ident;
    let generics = &input.generics;
    let mut body = TokenStream2::new();
    let mut encode_field_types = TokenStream2::new();
    match &input.fields {
        Fields::Named(fields) => {
            for field in &fields.named {
                let field_name = field.ident.as_ref().unwrap();
                let delta = quote! {
                    obi::OBIEncode::encode(&self.#field_name, writer)?;
                };
                body.extend(delta);

                let field_type = &field.ty;
                encode_field_types.extend(quote! {
                    #field_type: obi::enc::OBIEncode,
                });
            }
        }
        Fields::Unnamed(fields) => {
            for field_idx in 0..fields.unnamed.len() {
                let field_idx = Index {
                    index: field_idx as u32,
                    span: Span::call_site(),
                };
                let delta = quote! {
                    obi::OBIEncode::encode(&self.#field_idx, writer)?;
                };
                body.extend(delta);
            }
        }
        Fields::Unit => {}
    }
    Ok(quote! {
        impl #generics obi::enc::OBIEncode for #name #generics where #encode_field_types {
            fn encode<W: std::io::Write>(&self, writer: &mut W) -> std::result::Result<(), std::io::Error> {
                #body
                Ok(())
            }
        }
    })
}

// Rustfmt removes commas.
#[rustfmt::skip]
#[cfg(test)]
mod tests {
    use super::*;

    fn assert_eq(expected: TokenStream2, actual: TokenStream2) {
        assert_eq!(expected.to_string(), actual.to_string())
    }

    #[test]
    fn simple_struct() {
        let item_struct: ItemStruct = syn::parse2(quote!{
            struct A {
                x: u64,
                y: String,
            }
        }).unwrap();

        let actual = struct_enc(&item_struct).unwrap();
        let expected = quote!{
            impl obi::enc::OBIEncode for A
            where
                u64: obi::enc::OBIEncode,
                String: obi::enc::OBIEncode,
            {
                fn encode<W: std::io::Write>(&self, writer: &mut W) -> std::result::Result<(), std::io::Error> {
                    obi::OBIEncode::encode(&self.x, writer)?;
                    obi::OBIEncode::encode(&self.y, writer)?;
                    Ok(())
                }
            }
        };
        assert_eq(expected, actual);
    }
}
