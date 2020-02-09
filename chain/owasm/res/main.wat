(module
  (type $t0 (func))
  (type $t1 (func (result i64)))
  (type $t2 (func (param i32 i32)))
  (import "env" "getCurrentRequestID" (func $getCurrentRequestID (type $t1)))
  (import "env" "getRequestedValidatorCount" (func $getRequestedValidatorCount (type $t1)))
  (import "env" "saveReturnData" (func $saveReturnData (type $t2)))
  (func $execute (export "execute") (type $t0)
    (local $l0 i32)
    get_global $g0
    i32.const 16
    i32.sub
    tee_local $l0
    set_global $g0
    get_local $l0
    call $getCurrentRequestID
    call $getRequestedValidatorCount
    i64.add
    i64.const 42
    i64.add
    i64.store offset=8
    get_local $l0
    i32.const 8
    i32.add
    i32.const 8
    call $saveReturnData
    get_local $l0
    i32.const 16
    i32.add
    set_global $g0)
  (table $T0 1 1 anyfunc)
  (memory $memory (export "memory") 17)
  (global $g0 (mut i32) (i32.const 1049600)))
