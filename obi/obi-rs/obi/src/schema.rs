//! The important components are: `OBISchema` trait, `Definition` and `Declaration` types
//! * `OBISchema` trait allows any type that implements it to be self-descriptive, i.e. generate it's own schema;
//! * `Declaration` is used to describe the type identifier, e.g. `[u64]`;
//! * `Definition` is used to describe the structure of the type;

#![allow(dead_code)] // Unclear why rust check complains on fields of `Definition` variants.
use std::collections::hash_map::Entry;
use std::collections::*;

/// The type that we use to represent the declaration of the OBI type.
pub type Declaration = String;
/// The name of the field in the struct (can be used to convert JSON to OBI using the schema).
pub type FieldName = String;

/// The type that we use to represent the definition of the OBI type.
#[derive(PartialEq, Debug)]
pub enum Definition {
    /// A sequence of elements of length known at the run time and the same-type elements.
    Sequence { elements: Declaration },
    /// A structure, structurally similar to a tuple.
    Struct {
        fields: Vec<(FieldName, Declaration)>,
    },
}

pub fn get_schema(
    declaration: Declaration,
    definitions: &HashMap<Declaration, Definition>,
) -> String {
    match definitions.get(&declaration) {
        Some(definition) => match definition {
            Definition::Sequence { elements } => format!(
                "{}:[{}]",
                declaration,
                if elements == "u8" {
                    String::from("bytes")
                } else {
                    get_schema(elements.clone(), definitions)
                }
            ),
            Definition::Struct { fields } => format!(
                "{}}}",
                fields
                    .iter()
                    .enumerate()
                    .fold(String::from("{"), |acc, (i, (name, dec))| {
                        format!(
                            "{}{}{}:{}",
                            acc,
                            if i == 0 { "" } else { "," },
                            name,
                            get_schema(dec.clone(), definitions)
                        )
                    })
            ),
        },
        None => declaration,
    }
}

/// The declaration and the definition of the type that can be used to decode/encode OBI without
/// the Rust type that produced it.
pub trait OBISchema {
    /// Recursively, using DFS, add type definitions required for this type. For primitive types
    /// this is an empty map. Type definition explains how to serialize/deserialize a type.
    fn add_definitions_recursively(definitions: &mut HashMap<Declaration, Definition>);

    /// Helper method to add a single type definition to the map.
    fn add_definition(
        declaration: Declaration,
        definition: Definition,
        definitions: &mut HashMap<Declaration, Definition>,
    ) {
        match definitions.entry(declaration) {
            Entry::Occupied(occ) => {
                let existing_def = occ.get();
                assert_eq!(existing_def, &definition, "Redefining type schema for the same type name. Types with the same names are not supported.");
            }
            Entry::Vacant(vac) => {
                vac.insert(definition);
            }
        }
    }
    /// Get the name of the type without brackets.
    fn declaration() -> Declaration;
}

macro_rules! impl_for_renamed_primitives {
    ($($type: ident : $name: ident)+) => {
    $(
        impl OBISchema for $type {
            fn add_definitions_recursively(_definitions: &mut HashMap<Declaration, Definition>) {}
            fn declaration() -> Declaration {
                stringify!($name).to_string()
            }
        }
    )+
    };
}

macro_rules! impl_for_primitives {
    ($($type: ident)+) => {
    impl_for_renamed_primitives!{$($type : $type)+}
    };
}

impl_for_primitives!(bool char i8 i16 i32 i64 i128 u8 u16 u32 u64 u128);
impl_for_renamed_primitives!(String: string);

impl<T> OBISchema for Vec<T>
where
    T: OBISchema,
{
    fn add_definitions_recursively(definitions: &mut HashMap<Declaration, Definition>) {
        let definition = Definition::Sequence {
            elements: T::declaration(),
        };
        Self::add_definition(Self::declaration(), definition, definitions);
        T::add_definitions_recursively(definitions);
    }

    fn declaration() -> Declaration {
        format!(r#"[{}]"#, T::declaration())
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    macro_rules! map(
    () => { ::std::collections::HashMap::new() };
    { $($key:expr => $value:expr),+ } => {
        {
            let mut m = ::std::collections::HashMap::new();
            $(
                m.insert($key.to_string(), $value);
            )+
            m
        }
     };
    );

    #[test]
    fn simple_vec() {
        let actual_name = Vec::<u64>::declaration();
        let mut actual_defs = map!();
        Vec::<u64>::add_definitions_recursively(&mut actual_defs);
        assert_eq!("[u64]", actual_name);
        assert_eq!(
            map! {
            "[u64]" => Definition::Sequence { elements: "u64".to_string() }
            },
            actual_defs
        );
    }

    #[test]
    fn nested_vec() {
        let actual_name = Vec::<Vec<u64>>::declaration();
        let mut actual_defs = map!();
        Vec::<Vec<u64>>::add_definitions_recursively(&mut actual_defs);
        assert_eq!("[[u64]]", actual_name);
        assert_eq!(
            map! {
            "[u64]" => Definition::Sequence { elements: "u64".to_string() },
            "[[u64]]" => Definition::Sequence { elements: "[u64]".to_string() }
            },
            actual_defs
        );
    }
}
