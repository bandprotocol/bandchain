mod env;
mod error;
mod span;
mod vm;

use env::Env;
use error::Error;
use parity_wasm::builder;
use parity_wasm::elements::{self, MemorySection, MemoryType, Module};

use pwasm_utils::{self, rules};
use span::Span;
use std::ffi::c_void;
use wabt::wat2wasm;
use wasmer_runtime::{instantiate, Ctx};
use wasmer_runtime_core::error::RuntimeError;
use wasmer_runtime_core::{func, imports, wasmparser, Func};

// inspired by https://github.com/CosmWasm/cosmwasm/issues/81
// 512 pages = 32mb
static MEMORY_LIMIT: u32 = 512; // in pages

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
pub extern "C" fn do_run(code: Span, gas_limit: u32, is_prepare: bool, env: Env) -> Error {
    match run(code.read(), gas_limit, is_prepare, env) {
        Ok(_) => Error::NoError,
        Err(e) => e,
    }
}

#[no_mangle]
pub extern "C" fn do_wat2wasm(input: Span, output: &mut Span) -> Error {
    match wat2wasm(input.read()) {
        Ok(_wasm) => output.write(&_wasm),
        Err(e) => match e.kind() {
            wabt::ErrorKind::Parse(_) => Error::ParseError,
            wabt::ErrorKind::WriteBinary => Error::WriteBinaryError,
            wabt::ErrorKind::ResolveNames(_) => Error::ResolveNamesError,
            wabt::ErrorKind::Validate(_) => Error::ValidateError,
            _ => Error::UnknownError,
        },
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

fn inject_gas(module: Module) -> Result<Module, Error> {
    // Simple gas rule. Every opcode and memory growth costs 1 gas.
    let gas_rules = rules::Set::new(1, Default::default()).with_grow_cost(1);
    pwasm_utils::inject_gas_counter(module, &gas_rules).map_err(|_| Error::GasCounterInjectionError)
}

fn compile(code: &[u8]) -> Result<Vec<u8>, Error> {
    // Check that the given Wasm code is indeed a valid Wasm.
    wasmparser::validate(code, None).map_err(|_| Error::ValidateError)?;
    // Start the compiling chains. TODO: Add more safeguards.
    let module = elements::deserialize_buffer(code).map_err(|_| Error::DeserializationError)?;
    let module = inject_memory(module)?;
    let module = inject_gas(module)?;
    // Serialize the final Wasm code back to bytes.
    elements::serialize(module).map_err(|_| Error::SerializationError)
}

struct ImportReference(*mut c_void);
unsafe impl Send for ImportReference {}
unsafe impl Sync for ImportReference {}

fn run(code: &[u8], gas_limit: u32, is_prepare: bool, env: Env) -> Result<(), Error> {
    let vm = &mut vm::VMLogic::new(env, gas_limit);
    let raw_ptr = vm as *mut _ as *mut c_void;
    let import_reference = ImportReference(raw_ptr);
    let import_object = imports! {
        move || (import_reference.0, (|_: *mut c_void| {}) as fn(*mut c_void)),
        "env" => {
            "gas" => func!(|ctx: &mut Ctx, gas: u32| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.consume_gas(gas)
            }),
            "get_calldata_size" => func!(|ctx: &mut Ctx| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.get_calldata().len as i64
            }),
            "read_calldata" => func!(|ctx: &mut Ctx, ptr: i64, len: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                let calldata = vm.get_calldata();
                for (byte, cell) in calldata.read().iter().zip(ctx.memory(0).view()[ptr as usize..(ptr + len) as usize].iter()) { cell.set(*byte); }
            }),
            "set_return_data" => func!(|ctx: &mut Ctx, ptr: i64, len: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
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
                let data: Vec<u8> = ctx.memory(0).view()[ptr as usize..(ptr + len) as usize].iter().map(|cell| cell.get()).collect();
                vm.ask_external_data(eid, did, &data)
            }),
            "get_external_data_status" => func!(|ctx: &mut Ctx, eid: i64, vid: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.get_external_data_status(eid, vid)
            }),
            "get_external_data_size" => func!(|ctx: &mut Ctx, eid: i64, vid: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                vm.get_external_data(eid, vid).len as i64
            }),
            "read_external_data" => func!(|ctx: &mut Ctx, eid: i64, vid: i64, ptr: i64, len: i64| {
                let vm: &mut vm::VMLogic = unsafe { &mut *(ctx.data as *mut vm::VMLogic) };
                let calldata = vm.get_external_data(eid, vid);
                for (byte, cell) in calldata.read().iter().zip(ctx.memory(0).view()[ptr as usize..(ptr + len) as usize].iter()) { cell.set(*byte); }
            }),
        },
    };
    let instance = instantiate(code, &import_object).map_err(|_| Error::CompliationError)?;
    let entry = if is_prepare { "prepare" } else { "execute" };
    let function: Func<(), ()> = instance
        .exports
        .get(entry)
        .map_err(|_| Error::FunctionNotFoundError)?;
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
    use wabt::wat2wasm;

    fn get_module_from_wasm(code: &[u8]) -> Module {
        match elements::deserialize_buffer(code) {
            Ok(deserialized) => deserialized,
            Err(_) => panic!("Cannot deserialized"),
        }
    }

    #[test]
    fn test_inject_memory_ok() {
        let wasm = wat2wasm(r#"(module (memory 1))"#).unwrap();
        let module = get_module_from_wasm(&wasm);

        assert_matches!(inject_memory(module), Ok(_));
    }
    #[test]
    fn test_inject_memory_no_memory() {
        let wasm = wat2wasm("(module)").unwrap();
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
        let wasm_ok = wat2wasm("(module (memory 512))").unwrap();
        let module = get_module_from_wasm(&wasm_ok);
        assert_matches!(inject_memory(module), Ok(_));

        let wasm_too_big = wat2wasm("(module (memory 513))").unwrap();
        let module = get_module_from_wasm(&wasm_too_big);
        assert_eq!(inject_memory(module), Err(Error::MinimumMemoryExceedError));
    }

    #[test]
    fn test_inject_memory_maximum_size() {
        let wasm = wat2wasm("(module (memory 1 5))").unwrap();
        let module = get_module_from_wasm(&wasm);

        assert_eq!(inject_memory(module), Err(Error::SetMaximumMemoryError));
    }
}
