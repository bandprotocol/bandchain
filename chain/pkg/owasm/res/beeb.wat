(module
  (import "env" "requestExternalData" (func (param i64 i64 i32 i64) (result i64)))
  (import "env" "saveReturnData" (func (param i32 i64) (result i64)))
  (func
    i64.const 1
    i64.const 1
    i32.const 1048576
    i64.const 4
    call 0
    drop
    i64.const 2
    i64.const 2
    i32.const 1048576
    i64.const 4
    call 0
    drop
    i64.const 3
    i64.const 3
    i32.const 1048576
    i64.const 4
    call 0
    drop)
  (func
    i32.const 1048576
    i64.const 4
    call 1
    drop)
  (memory 17)
  (data (i32.const 1048576) "beeb")
  (export "prepare" (func 2))
  (export "execute" (func 3)))
