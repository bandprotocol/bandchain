use proc_macro2::TokenStream;
use quote::quote;
use syn::Generics;

pub fn declaration(ident_str: &String, generics: &Generics) -> (TokenStream, Vec<TokenStream>) {
    let (_, _, where_clause_generics) = generics.split_for_impl();
    // Generate function that returns the name of the type.
    let mut declaration_params = vec![];
    let mut where_clause = vec![];
    if let Some(where_clause_generics) = where_clause_generics {
        let where_clause_generics = &where_clause_generics.predicates;
        where_clause.push(quote! {#where_clause_generics});
    }
    for type_param in generics.type_params() {
        let type_param_name = &type_param.ident;
        declaration_params.push(quote! {
            <#type_param_name>::declaration()
        });
        where_clause.push(quote! {
            #type_param_name: obi::OBISchema
        });
    }
    let result = if declaration_params.is_empty() {
        quote! {
                #ident_str.to_string()
        }
    } else {
        quote! {
                let params = vec![#(#declaration_params),*];
                format!(r#"{}<{}>"#, #ident_str, params.join(", "))
        }
    };
    (result, where_clause)
}
