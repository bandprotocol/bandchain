use std::env;
use std::fs;
use wasmer_runtime::{error, imports, instantiate, Instance, Value};

fn store_bytes(instance: &Instance, data: &[u8]) -> Option<u64> {
    let sz = data.len();
    let loc = instance.call("__allocate", &[Value::I32(sz as i32)]).ok()?[0].to_u128() as u32;
    let memory = instance.context().memory(0);
    for (idx, ch) in data.iter().enumerate() {
        memory.view()[(loc + (idx as u32)) as usize].set(*ch);
    }
    Some(((sz as u64) << 32) | ((loc as u32) as u64))
}

fn main() -> error::Result<()> {
    let args = env::args().collect::<Vec<String>>();
    let code = fs::read(&args[1]).unwrap();
    let param_json = args[2].as_bytes();

    let instance = instantiate(&code, &imports! {})?;
    let params_input = store_bytes(&instance, param_json).unwrap();
    let ptr_and_sz = instance.call("__serialize_params", &[Value::I64(params_input as i64)])?[0]
        .to_u128() as i64;
    let memory = instance.context().memory(0);
    let ptr = (ptr_and_sz & ((1 << 32) - 1)) as usize;
    let sz = (ptr_and_sz >> 32) as usize;
    let params_hex: Vec<u8> =
        memory.view()[ptr..(ptr + sz)].iter().map(|cell| cell.get()).collect();
    println!(
        "{}",
        params_hex.iter().map(|b| format!("{:02x}", b)).collect::<Vec<String>>().join("")
    );
    Ok(())
}
