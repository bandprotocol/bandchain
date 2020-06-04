(module
  (type $t0 (func (result i64)))
  (import "env" "get_ask_count" (func $getAskCount (type $t0)))
  (func $add_one (export "prepare") (type $t0) (result i64)
    call $getAskCount
    i64.const 10
    i64.add)
  (table $T0 1 1 anyfunc)
  (memory $memory (export "memory") 17))
