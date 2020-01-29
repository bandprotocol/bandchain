use owasm::core::{decode_cmds, encode_outputs, execute_with_local_env};
use std::env;
use std::fs;
use wasmer_runtime::{error, imports, instantiate, Instance, Value};

fn decode_hex(s: &str) -> Option<Vec<u8>> {
    (0..s.len()).step_by(2).map(|i| u8::from_str_radix(&s[i..i + 2], 16).ok()).collect()
}

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
    let code = fs::read(&env::args().collect::<Vec<String>>()[1]).unwrap();
    let param_raw = decode_hex(&env::args().collect::<Vec<String>>()[2]).unwrap();
    let input_raw: Vec<u8>;
    let output_raw: Vec<u8>;
    {
        let instance = instantiate(&code, &imports! {})?;
        let params_input = store_bytes(&instance, &param_raw).unwrap();
        let ptr_and_sz =
            instance.call("__prepare", &[Value::I64(params_input as i64)])?[0].to_u128() as i64;
        let memory = instance.context().memory(0);
        let ptr = (ptr_and_sz & ((1 << 32) - 1)) as usize;
        let sz = (ptr_and_sz >> 32) as usize;
        input_raw = memory.view()[ptr..(ptr + sz)].iter().map(|cell| cell.get()).collect();
    }
    {
        let output =
            vec![encode_outputs(execute_with_local_env(decode_cmds(&input_raw).unwrap())).unwrap()];
        let instance = instantiate(&code, &imports! {})?;
        // Allocate and fill parameter
        let params_input = store_bytes(&instance, &param_raw).unwrap();

        // Allocate and fill data
        let memory = instance.context().memory(0);
        let list_ptr = instance.call("__allocate", &[Value::I32(output.len() as i32 * 8)])?[0]
            .to_u128() as usize;

        for (idx_order, each) in output.iter().enumerate() {
            let data_ptr = store_bytes(&instance, &each).unwrap().to_le_bytes();
            for (idx, ch) in data_ptr.iter().enumerate() {
                memory.view()[(list_ptr + idx_order * 8 + idx)].set(*ch);
            }
        }

        let list_ptr_and_size = ((output.len() as u64) << 32) | (list_ptr as u64);

        let out_ptr_and_size = instance.call(
            "__execute",
            &[Value::I64(params_input as i64), Value::I64(list_ptr_and_size as i64)],
        )?[0]
            .to_u128() as u64;
        let ptr = (out_ptr_and_size & ((1 << 32) - 1)) as usize;
        let sz = (out_ptr_and_size >> 32) as usize;
        output_raw = memory.view()[ptr..(ptr + sz)].iter().map(|cell| cell.get()).collect();
    }
    println!(
        "0x{}",
        output_raw.iter().map(|b| format!("{:02x}", b)).collect::<Vec<String>>().join("")
    );
    Ok(())
}
