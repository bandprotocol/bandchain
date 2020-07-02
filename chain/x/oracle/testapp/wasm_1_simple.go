package testapp

// A simple Owasm script with the following specification:
//   PREPARE:
//     CALL ask_external_data with EID 1 DID 1 CALLDATA "beeb"
//     CALL ask_external_data with EID 2 DID 2 CALLDATA "beeb"
//     CALL ask_external_data with EID 3 DID 3 CALLDATA "beeb"
//   EXECUTE:
//     CALL set_return_date with RETDATE "beeb"
var Wasm1 []byte = wat2wasm([]byte(`
(module
	(type $t0 (func))
	(type $t1 (func (param i64 i64 i64 i64)))
	(type $t2 (func (param i64 i64)))
	(import "env" "ask_external_data" (func $ask_external_data (type $t1)))
	(import "env" "set_return_data" (func $set_return_data (type $t2)))
	(func $prepare (export "prepare") (type $t0)
	  (local $l0 i64)
	  i64.const 1
	  i64.const 1
	  i32.const 1024
	  i64.extend_u/i32
	  tee_local $l0
	  i64.const 4
	  call $ask_external_data
	  i64.const 2
	  i64.const 2
	  get_local $l0
	  i64.const 4
	  call $ask_external_data
	  i64.const 3
	  i64.const 3
	  get_local $l0
	  i64.const 4
	  call $ask_external_data)
	(func $execute (export "execute") (type $t0)
	  i32.const 1024
	  i64.extend_u/i32
	  i64.const 4
	  call $set_return_data)
	(table $T0 1 1 anyfunc)
	(memory $memory (export "memory") 17)
	(data (i32.const 1024) "beeb"))
`))
