use crate::helpers::declaration;
use proc_macro2::TokenStream;
use quote::{quote, ToTokens};
use syn::{Fields, ItemStruct};

pub fn process_struct(input: &ItemStruct) -> syn::Result<TokenStream> {
    let name = &input.ident;
    let name_str = name.to_token_stream().to_string();
    let generics = &input.generics;
    let (impl_generics, ty_generics, _) = generics.split_for_impl();
    // Generate function that returns the name of the type.
    let (declaration, mut where_clause) = declaration(&name_str, &input.generics);

    // Generate function that returns the schema of required types.
    let mut fields_vec = vec![];
    let mut struct_fields = TokenStream::new();
    let mut add_definitions_recursively_rec = TokenStream::new();
    match &input.fields {
        Fields::Named(fields) => {
            for field in &fields.named {
                let field_name = field.ident.as_ref().unwrap().to_token_stream().to_string();
                let field_type = &field.ty;
                fields_vec.push(quote! {
                    (#field_name.to_string(), <#field_type>::declaration())
                });
                add_definitions_recursively_rec.extend(quote! {
                    <#field_type>::add_definitions_recursively(definitions);
                });
                where_clause.push(quote! {
                    #field_type: obi::OBISchema
                });
            }
            if !fields_vec.is_empty() {
                struct_fields = quote! {
                    let fields = vec![#(#fields_vec),*];
                };
            }
        }
        // Unsupported on unnamed struct (tuple)
        Fields::Unnamed(_fields) => {}
        Fields::Unit => {}
    }

    if fields_vec.is_empty() {
        struct_fields = quote! {
            let fields = vec![];
        };
    }

    let add_definitions_recursively = quote! {
        fn add_definitions_recursively(definitions: &mut ::std::collections::HashMap<obi::schema::Declaration, obi::schema::Definition>) {
            #struct_fields
            let definition = obi::schema::Definition::Struct { fields };
            Self::add_definition(Self::declaration(), definition, definitions);
            #add_definitions_recursively_rec
        }
    };
    let where_clause = if !where_clause.is_empty() {
        quote! { where #(#where_clause),*}
    } else {
        TokenStream::new()
    };
    Ok(quote! {
        impl #impl_generics obi::OBISchema for #name #ty_generics #where_clause {
            fn declaration() -> obi::schema::Declaration {
                #declaration
            }
            #add_definitions_recursively
        }
    })
}

// Rustfmt removes comas.
#[rustfmt::skip::macros(quote)]
#[cfg(test)]
mod tests {
    use super::*;

    fn assert_eq(expected: TokenStream, actual: TokenStream) {
        assert_eq!(expected.to_string(), actual.to_string())
    }

    #[test]
    fn unit_struct() {
        let item_struct: ItemStruct = syn::parse2(quote!{
            struct A;
        })
        .unwrap();

        let actual = process_struct(&item_struct).unwrap();
        let expected = quote!{
            impl obi::OBISchema for A
            {
                fn declaration() -> obi::schema::Declaration {
                    "A".to_string()
                }
                fn add_definitions_recursively(definitions: &mut ::std::collections::HashMap<obi::schema::Declaration, obi::schema::Definition>) {
                    let fields = vec![];
                    let definition = obi::schema::Definition::Struct { fields };
                    Self::add_definition(Self::declaration(), definition, definitions);
                }
            }
        };
        assert_eq(expected, actual);
    }

    #[test]
    fn simple_struct() {
        let item_struct: ItemStruct = syn::parse2(quote!{
            struct A {
                x: u64,
                y: String,
            }
        })
        .unwrap();

        let actual = process_struct(&item_struct).unwrap();
        let expected = quote!{
            impl obi::OBISchema for A
            where
                u64: obi::OBISchema,
                String: obi::OBISchema
            {
                fn declaration() -> obi::schema::Declaration {
                    "A".to_string()
                }
                fn add_definitions_recursively(
                    definitions: &mut ::std::collections::HashMap<
                        obi::schema::Declaration,
                        obi::schema::Definition
                    >
                ) {
                    let fields = vec![
                        ("x".to_string(), <u64>::declaration()),
                        ("y".to_string(), <String>::declaration())
                    ];
                    let definition = obi::schema::Definition::Struct { fields };
                    Self::add_definition(Self::declaration(), definition, definitions);
                    <u64>::add_definitions_recursively(definitions);
                    <String>::add_definitions_recursively(definitions);
                }
            }
        };
        assert_eq(expected, actual);
    }

    #[test]
    fn simple_generics() {
        let item_struct: ItemStruct = syn::parse2(quote!{
            struct A<K, V> {
                x: HashMap<K, V>,
                y: String,
            }
        })
        .unwrap();

        let actual = process_struct(&item_struct).unwrap();
        let expected = quote!{
            impl<K, V> obi::OBISchema for A<K, V>
            where
                K: obi::OBISchema,
                V: obi::OBISchema,
                HashMap<K, V>: obi::OBISchema,
                String: obi::OBISchema
            {
                fn declaration() -> obi::schema::Declaration {
                    let params = vec![<K>::declaration(), <V>::declaration()];
                    format!(r#"{}<{}>"#, "A", params.join(", "))
                }
                fn add_definitions_recursively(
                    definitions: &mut ::std::collections::HashMap<
                        obi::schema::Declaration,
                        obi::schema::Definition
                    >
                ) {
                    let fields = vec![
                        ("x".to_string(), <HashMap<K, V> >::declaration()),
                        ("y".to_string(), <String>::declaration())
                    ];
                    let definition = obi::schema::Definition::Struct { fields };
                    Self::add_definition(Self::declaration(), definition, definitions);
                    <HashMap<K, V> >::add_definitions_recursively(definitions);
                    <String>::add_definitions_recursively(definitions);
                }
            }
        };
        assert_eq(expected, actual);
    }

    #[test]
    fn tuple_struct_whole_skip() {
        let item_struct: ItemStruct = syn::parse2(quote!{
            struct A(#[obi_skip] String);
        })
        .unwrap();

        let actual = process_struct(&item_struct).unwrap();
        let expected = quote!{
            impl obi::OBISchema for A {
                fn declaration() -> obi::schema::Declaration {
                    "A".to_string()
                }
                fn add_definitions_recursively(
                    definitions: &mut ::std::collections::HashMap<
                        obi::schema::Declaration,
                        obi::schema::Definition
                    >
                ) {
                    let fields = vec![];
                    let definition = obi::schema::Definition::Struct { fields };
                    Self::add_definition(Self::declaration(), definition, definitions);
                }
            }
        };
        assert_eq(expected, actual);
    }
}
