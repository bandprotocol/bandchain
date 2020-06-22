mod env;
mod error;
mod span;
mod vm;

use env::Env;
use error::Error;
use parity_wasm::builder;
use parity_wasm::elements::{self, External, ImportEntry, MemoryType, Module};

use pwasm_utils::{self, rules};
use span::Span;
use std::ffi::c_void;
use wasmer_runtime::{instantiate, Ctx};
use wasmer_runtime_core::error::RuntimeError;
use wasmer_runtime_core::{func, imports, wasmparser, Func};

// inspired by https://github.com/CosmWasm/cosmwasm/issues/81
// 512 pages = 32mb
static MEMORY_LIMIT: u32 = 512; // in pages
static MAX_STACK_HEIGHT: u32 = 16 * 1024; // 16Kib of stack.

static REQUIRED_EXPORTS: &[&str] = &["prepare", "execute"];
static SUPPORTED_IMPORTS: &[&str] = &[
    "env.get_calldata_size",
    "env.read_calldata",
    "env.set_return_data",
    "env.get_ask_count",
    "env.get_min_count",
    "env.get_ans_count",
    "env.ask_external_data",
    "env.get_external_data_status",
    "env.get_external_data_size",
    "env.read_external_data",
];

#[no_mangle]
pub extern "C" fn do_compile(input: Span, output: &mut Span) -> Error {
    match compile(input.read()) {
        Ok(out) => {
            output.write(&out);
            Error::NoError
        }
        Err(e) => e,
    }
}

#[no_mangle]
pub extern "C" fn do_run(
    code: Span,
    gas_limit: u32,
    span_size: i64,
    is_prepare: bool,
    env: Env,
) -> Error {
    match run(code.read(), gas_limit, span_size, is_prepare, env) {
        Ok(_) => Error::NoError,
        Err(e) => e,
    }
}

fn inject_memory(module: Module) -> Result<Module, Error> {
    let mut m = module;
    let section = match m.memory_section() {
        Some(section) => section,
        None => return Err(Error::NoMemoryWasmError),
    };

    // The valid wasm has only the section length of memory.
    // We check the wasm is valid in the first step of compile fn.
    let memory = section.entries()[0];
    let limits = memory.limits();

    if limits.initial() > MEMORY_LIMIT {
        return Err(Error::MinimumMemoryExceedError);
    }

    if limits.maximum() != None {
        return Err(Error::SetMaximumMemoryError);
    }

    // set max memory page = MEMORY_LIMIT
    let memory = MemoryType::new(limits.initial(), Some(MEMORY_LIMIT));

    // Memory existance already checked
    let entries = m.memory_section_mut().unwrap().entries_mut();
    entries.pop();
    entries.push(memory);

    Ok(builder::from_module(m).build())
}

fn inject_stack_height(module: Module) -> Result<Module, Error> {
    pwasm_utils::stack_height::inject_limiter(module, MAX_STACK_HEIGHT)
        .map_err(|_| Error::StackHeightInstrumentationError)
}

fn inject_gas(module: Module) -> Result<Module, Error> {
    // Simple gas rule. Every opcode and memory growth costs 1 gas.
    let gas_rules = rules::Set::new(1, Default::default()).with_grow_cost(1);
    pwasm_utils::inject_gas_counter(module, &gas_rules).map_err(|_| Error::GasCounterInjectionError)
}

fn check_wasm_exports(module: &Module) -> Result<(), Error> {
    let available_exports: Vec<&str> = module.export_section().map_or(vec![], |export_section| {
        export_section
            .entries()
            .iter()
            .map(|entry| entry.field())
            .collect()
    });

    for required_export in REQUIRED_EXPORTS {
        if !available_exports.contains(required_export) {
            return Err(Error::CheckWasmExportsError);
        }
    }

    Ok(())
}

fn check_wasm_imports(module: &Module) -> Result<(), Error> {
    let required_imports: Vec<ImportEntry> = module
        .import_section()
        .map_or(vec![], |import_section| import_section.entries().to_vec());

    for required_import in required_imports {
        let full_name = format!("{}.{}", required_import.module(), required_import.field());
        if !SUPPORTED_IMPORTS.contains(&full_name.as_str()) {
            return Err(Error::CheckWasmImportsError);
        }

        match required_import.external() {
            External::Function(_) => {} // ok
            _ => return Err(Error::CheckWasmImportsError),
        };
    }
    Ok(())
}

fn compile(code: &[u8]) -> Result<Vec<u8>, Error> {
    // Check that the given Wasm code is indeed a valid Wasm.
    wasmparser::validate(code, None).map_err(|_| Error::ValidateError)?;
    // Start the compiling chains. TODO: Add more safeguards.
    let module = elements::deserialize_buffer(code).map_err(|_| Error::DeserializationError)?;
    check_wasm_exports(&module)?;
    check_wasm_imports(&module)?;
    let module = inject_memory(module)?;
    let module = inject_gas(module)?;
    let module = inject_stack_height(module)?;
    // Serialize the final Wasm code back to bytes.
    elements::serialize(module).map_err(|_| Error::SerializationError)
}

struct ImportReference(*mut c_void);
unsafe impl Send for ImportReference {}
unsafe impl Sync for ImportReference {}

fn run(
    code: &[u8],
    gas_limit: u32,
    span_size: i64,
    is_prepare: bool,
    env: Env,
) -> Result<(), Error> {
    let vm = &mut vm::VMLogic::new(env, gas_limit, span_size);
    let raw_ptr = vm as *mut _ as *mut c_void;
    let import_reference = ImportReference(raw_ptr);
    let import_object = imports! {
        move || (import_reference.0, (|_: *mut c_void| {}) as fn(*mut c_void)),
        "env" => {
            "gas" => func!(|ctx: &mut Ctx, gas: u32| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.consume_gas(gas)
            }),
            // TODO: Change specification of OEI
            "get_calldata_size" => func!(|ctx: &mut Ctx| -> Result<i64, Error> {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                let span_size = vm.get_span_size() as usize;
                let mut mem: Vec<u8> = Vec::with_capacity(span_size);
                let mut calldata = Span::create_writable(mem.as_mut_ptr(), span_size);
                vm.get_calldata(&mut calldata)?;
                Ok(calldata.len as i64)
            }),
            "read_calldata" => func!(|ctx: &mut Ctx, ptr: i64, len: i64| -> Result<(), Error>{
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                let span_size = vm.get_span_size() as usize;
                let mut mem: Vec<u8> = Vec::with_capacity(span_size);
                let mut calldata = Span::create_writable(mem.as_mut_ptr(), span_size);
                vm.get_calldata(&mut calldata)?;
                for (byte, cell) in calldata.read().iter().zip(ctx.memory(0).view()[ptr as usize..(ptr + len) as usize].iter()) { cell.set(*byte); }
                Ok(())
            }),
            "set_return_data" => func!(|ctx: &mut Ctx, ptr: i64, len: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                if len as usize > vm.get_span_size() {
                    return Err(Error::SpanExceededCapacityError);
                }
                let data: Vec<u8> = ctx.memory(0).view()[ptr as usize..(ptr + len) as usize].iter().map(|cell| cell.get()).collect();
                vm.set_return_data(&data)
            }),
            "get_ask_count" => func!(|ctx: &mut Ctx| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.get_ask_count()
            }),
            "get_min_count" => func!(|ctx: &mut Ctx| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.get_min_count()
            }),
            "get_ans_count" => func!(|ctx: &mut Ctx| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.get_ans_count()
            }),
            "ask_external_data" => func!(|ctx: &mut Ctx, eid: i64, did: i64, ptr: i64, len: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                if len as usize > vm.get_span_size() {
                    return Err(Error::SpanExceededCapacityError);
                }
                let data: Vec<u8> = ctx.memory(0).view()[ptr as usize..(ptr + len) as usize].iter().map(|cell| cell.get()).collect();
                vm.ask_external_data(eid, did, &data)
            }),
            "get_external_data_status" => func!(|ctx: &mut Ctx, eid: i64, vid: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.get_external_data_status(eid, vid)
            }),
            // TODO: Change specification of OEI
            "get_external_data_size" => func!(|ctx: &mut Ctx, eid: i64, vid: i64| -> Result<i64, Error>{
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                let span_size = vm.get_span_size() as usize;
                let mut mem: Vec<u8> = Vec::with_capacity(span_size);
                let mut data = Span::create_writable(mem.as_mut_ptr(), span_size);
                vm.get_external_data(eid, vid, &mut data)?;
                Ok(data.len as i64)
            }),
            "read_external_data" => func!(|ctx: &mut Ctx, eid: i64, vid: i64, ptr: i64, len: i64| -> Result<(),Error> {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                let span_size = vm.get_span_size() as usize;
                let mut mem: Vec<u8> = Vec::with_capacity(span_size);
                let mut data = Span::create_writable(mem.as_mut_ptr(), span_size);
                vm.get_external_data(eid, vid, &mut data)?;
                for (byte, cell) in data.read().iter().zip(ctx.memory(0).view()[ptr as usize..(ptr + len) as usize].iter()) { cell.set(*byte); }
                Ok(())
            }),
        },
    };
    let instance = instantiate(code, &import_object).map_err(|_| Error::CompliationError)?;
    let entry = if is_prepare { "prepare" } else { "execute" };
    let function: Func<(), ()> = instance
        .exports
        .get(entry)
        .map_err(|_| Error::InvalidSignatureFunctionError)?;
    function.call().map_err(|err| match err {
        RuntimeError::User(uerr) => {
            if let Some(err) = uerr.downcast_ref::<Error>() {
                err.clone()
            } else {
                Error::UnknownError
            }
        }
        _ => Error::RunError,
    })
}
#[cfg(test)]
mod test {
    use super::*;
    use assert_matches::assert_matches;
    use parity_wasm::elements;
    use std::io::{Read, Write};
    use std::process::Command;
    use tempfile::NamedTempFile;

    fn wat2wasm(wat: impl AsRef<[u8]>) -> Vec<u8> {
        let mut input_file = NamedTempFile::new().unwrap();
        let mut output_file = NamedTempFile::new().unwrap();
        input_file.write_all(wat.as_ref()).unwrap();
        Command::new("wat2wasm")
            .args(&[
                input_file.path().to_str().unwrap(),
                "-o",
                output_file.path().to_str().unwrap(),
            ])
            .output()
            .unwrap();
        let mut wasm = Vec::new();
        output_file.read_to_end(&mut wasm).unwrap();
        wasm
    }

    fn get_module_from_wasm(code: &[u8]) -> Module {
        match elements::deserialize_buffer(code) {
            Ok(deserialized) => deserialized,
            Err(_) => panic!("Cannot deserialized"),
        }
    }

    #[test]
    fn test_inject_memory_ok() {
        let wasm = wat2wasm(r#"(module (memory 1))"#);
        let module = get_module_from_wasm(&wasm);

        assert_matches!(inject_memory(module), Ok(_));
    }

    #[test]
    fn test_compile() {
        let wasm = wat2wasm(
            r#"(module
            (type (func (param i64 i64 i32 i64) (result i64)))
            (import "env" "ask_external_data" (func (type 0)))
            (func
              (local $idx i32)
              (set_local $idx (i32.const 0))
              (block
                  (loop
                    (set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
                    (br_if 0 (i32.lt_u (get_local $idx) (i32.const 1000000000)))
                  )
                )
            )
            (func (;"execute": Resolves with result "beeb";)
            )
            (memory 17)
            (data (i32.const 1048576) "beeb") (;str = "beeb";)
            (export "prepare" (func 0))
            (export "execute" (func 1)))
          "#,
        );
        let code = compile(&wasm).unwrap();

        let expected = wat2wasm(
            r#"(module
                (type (;0;) (func (param i64 i64 i32 i64) (result i64)))
                (type (;1;) (func))
                (type (;2;) (func (param i32)))
                (import "env" "ask_external_data" (func (;0;) (type 0)))
                (import "env" "gas" (func (;1;) (type 2)))
                (func (;2;) (type 1)
                  (local i32)
                  i32.const 4
                  call 1
                  i32.const 0
                  local.set 0
                  block  ;; label = @1
                    loop  ;; label = @2
                      i32.const 8
                      call 1
                      local.get 0
                      i32.const 1
                      i32.add
                      local.set 0
                      local.get 0
                      i32.const 1000000000
                      i32.lt_u
                      br_if 0 (;@2;)
                    end
                  end)
                (func (;3;) (type 1))
                (func (;4;) (type 1)
                  global.get 0
                  i32.const 3
                  i32.add
                  global.set 0
                  global.get 0
                  i32.const 16384
                  i32.gt_u
                  if  ;; label = @1
                    unreachable
                  end
                  call 2
                  global.get 0
                  i32.const 3
                  i32.sub
                  global.set 0)
                (memory (;0;) 17 512)
                (global (;0;) (mut i32) (i32.const 0))
                (export "prepare" (func 0))
                (export "execute" (func 4))
                (data (;0;) (i32.const 1048576) "beeb"))
          "#,
        );
        assert_eq!(code, expected);
    }

    #[test]
    fn test_inject_memory_no_memory() {
        let wasm = wat2wasm("(module)");
        let module = get_module_from_wasm(&wasm);

        assert_eq!(inject_memory(module), Err(Error::NoMemoryWasmError));
    }
    #[test]
    fn test_inject_memory_two_memories() {
        // Generated manually because wat2wasm protects us from creating such Wasm:
        // "error: only one memory block allowed"
        let wasm = hex::decode(concat!(
            "0061736d", // magic bytes
            "01000000", // binary version (uint32)
            "05",       // section type (memory)
            "05",       // section length
            "02",       // number of memories
            "0009",     // element of type "resizable_limits", min=9, max=unset
            "0009",     // element of type "resizable_limits", min=9, max=unset
        ))
        .unwrap();
        let r = compile(&wasm);
        assert_eq!(r, Err(Error::ValidateError));
    }

    #[test]
    fn test_inject_memory_initial_size() {
        let wasm_ok = wat2wasm("(module (memory 512))");
        let module = get_module_from_wasm(&wasm_ok);
        assert_matches!(inject_memory(module), Ok(_));

        let wasm_too_big = wat2wasm("(module (memory 513))");
        let module = get_module_from_wasm(&wasm_too_big);
        assert_eq!(inject_memory(module), Err(Error::MinimumMemoryExceedError));
    }

    #[test]
    fn test_inject_memory_maximum_size() {
        let wasm = wat2wasm("(module (memory 1 5))");
        let module = get_module_from_wasm(&wasm);

        assert_eq!(inject_memory(module), Err(Error::SetMaximumMemoryError));
    }

    #[test]
    fn test_inject_stack_height() {
        let wasm = wat2wasm(
            r#"(module
            (func
              (local $idx i32)
              (set_local $idx (i32.const 0))
              (block
                  (loop
                    (set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
                    (br_if 0 (i32.lt_u (get_local $idx) (i32.const 1000000000)))
                  )
                )
            )
            (func (;"execute": Resolves with result "beeb";)
            )
            (memory 17)
            (data (i32.const 1048576) "beeb") (;str = "beeb";)
            (export "prepare" (func 0))
            (export "execute" (func 1)))
          "#,
        );

        let module = inject_stack_height(get_module_from_wasm(&wasm)).unwrap();
        let wasm = elements::serialize(module).unwrap();

        let expected = wat2wasm(
            r#"(module
                (type (;0;) (func))
                (func (;0;) (type 0)
                  (local i32)
                  i32.const 0
                  local.set 0
                  block  ;; label = @1
                    loop  ;; label = @2
                      local.get 0
                      i32.const 1
                      i32.add
                      local.set 0
                      local.get 0
                      i32.const 1000000000
                      i32.lt_u
                      br_if 0 (;@2;)
                    end
                  end)
                (func (;1;) (type 0))
                (func (;2;) (type 0)
                  global.get 0
                  i32.const 3
                  i32.add
                  global.set 0
                  global.get 0
                  i32.const 16384
                  i32.gt_u
                  if  ;; label = @1
                    unreachable
                  end
                  call 0
                  global.get 0
                  i32.const 3
                  i32.sub
                  global.set 0)
                (memory (;0;) 17)
                (global (;0;) (mut i32) (i32.const 0))
                (export "prepare" (func 2))
                (export "execute" (func 1))
                (data (;0;) (i32.const 1048576) "beeb"))
          "#,
        );
        assert_eq!(wasm, expected);
    }

    #[test]
    fn test_check_wasm_imports() {
        let wasm = wat2wasm(
            r#"(module
                (type (func (param i64 i64 i32 i64) (result i64)))
                (import "env" "beeb" (func (type 0))))"#,
        );
        let module = get_module_from_wasm(&wasm);

        assert_eq!(
            check_wasm_imports(&module),
            Err(Error::CheckWasmImportsError)
        );

        let wasm = wat2wasm(
            r#"(module
                (type (func (param i64 i64 i32 i64) (result i64)))
                (import "env" "ask_external_data" (func  (type 0))))"#,
        );
        let module = get_module_from_wasm(&wasm);

        assert_eq!(check_wasm_imports(&module), Ok(()));
    }

    #[test]
    fn test_check_wasm_exports() {
        let wasm = wat2wasm(
            r#"(module
            (func $execute (export "execute")))"#,
        );
        let module = get_module_from_wasm(&wasm);

        assert_eq!(
            check_wasm_exports(&module),
            Err(Error::CheckWasmExportsError)
        );

        let wasm = wat2wasm(
            r#"(module
                (func $execute (export "execute"))
                (func $prepare (export "prepare"))
              )"#,
        );
        let module = get_module_from_wasm(&wasm);

        assert_eq!(check_wasm_exports(&module), Ok(()));
    }
}
