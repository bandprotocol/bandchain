use std::mem;

pub fn allocate(size: usize) -> *mut u8 {
    let mut buffer = Vec::with_capacity(size);
    let pointer = buffer.as_mut_ptr();
    mem::forget(buffer);
    pointer
}

pub fn allocate_const(output: &[u8]) -> *const u8 {
    let sz = output.len();
    let loc = allocate(4 + sz as usize);
    unsafe {
        std::ptr::copy_nonoverlapping(&(sz as u32), loc as *mut u32, 4);
        std::ptr::copy_nonoverlapping(output.as_ptr(), loc.offset(4), sz)
    };
    loc
}
