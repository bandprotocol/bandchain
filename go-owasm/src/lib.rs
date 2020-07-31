mod env;
mod span;
mod vm;

use env::{Env, RunOutput};
use owasm::core;
use span::Span;

#[no_mangle]
pub extern "C" fn do_compile(input: Span, output: &mut Span) -> owasm::core::error::Error {
    match core::compile(input.read()) {
        Ok(out) => {
            output.write(&out);
            owasm::core::error::Error::NoError
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
    output: &mut RunOutput,
) -> owasm::core::error::Error {
    let vm_env = vm::VMEnv::new(env, span_size);
    match core::run(code.read(), gas_limit, is_prepare, &vm_env) {
        Ok(gas_used) => {
            output.gas_used = gas_used;
            owasm::core::error::Error::NoError
        }
        Err(e) => e,
    }
}
