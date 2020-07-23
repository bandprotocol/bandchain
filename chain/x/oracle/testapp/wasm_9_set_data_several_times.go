package testapp

var Wasm9 []byte = wat2wasm([]byte(`
(module
	(type $t0 (func))
	(type $t1 (func (param i64 i64 i64 i64)))
	(type $t2 (func (param i64 i64)))
	(import "env" "ask_external_data" (func $ask_external_data (type $t1)))
	(import "env" "set_return_data" (func $set_return_data (type $t2)))
	(func $prepare (export "prepare") (type $t0))
	(func $execute (export "execute") (type $t0)
	  i32.const 1024
	  i64.extend_u/i32
	  i64.const 4
	  call $set_return_data
	  i32.const 1024
	  i64.extend_u/i32
	  i64.const 4
	  call $set_return_data)
	(table $T0 1 1 anyfunc)
	(memory $memory (export "memory") 17)
	(data (i32.const 1024) "beeb"))
`))
