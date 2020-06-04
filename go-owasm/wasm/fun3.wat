(module
  (type $t0 (func))
  (type $t1 (func (param i64 i64 i64 i64)))
  (import "env" "ask_external_data" (func $ask_external_data (type $t1)))
  (func $prepare (export "prepare") (type $t0)
    (local $l0 i64)
    i64.const 1
    i64.const 1
    i32.const 1024
    i64.extend_u/i32
    tee_local $l0
    i64.const 10
    call $ask_external_data
    i64.const 2
    i64.const 3
    get_local $l0
    i64.const 10
    call $ask_external_data
    i64.const 4
    i64.const 5
    get_local $l0
    i64.const 10
    call $ask_external_data)
  (table $T0 1 1 anyfunc)
  (memory $memory (export "memory") 17)
  (data (i32.const 1024) "helloworld"))
