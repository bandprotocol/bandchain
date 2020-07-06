package testapp

import (
	"fmt"
)

// An oracle script to test heavy computation.
//   PREPARE:
//     Loop for `time` times and ask external data with calldata "new beeb".
//   EXECUTE:
//     Loop for `time` times and set return data "new beeb".
func Wasm56(time int) []byte {
	return wat2wasm([]byte(fmt.Sprintf(`(module
	 (type $t0 (func))
	 (type $t1 (func (param i64 i64 i64 i64)))
	 (type $t2 (func (param i64 i64)))
	 (import "env" "ask_external_data" (func $ask_external_data (type $t1)))
	 (import "env" "set_return_data" (func $set_return_data (type $t2)))
	 (func $prepare (export "prepare") (type $t0)
	   (local $l0 i64)
	   (local $idx i32)
	   (set_local $idx (i32.const 0))
	   (block
		(loop
		(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
		(br_if 0 (i32.lt_u (get_local $idx) (i32.const %d)))
		)
	   )
	   i64.const 1
	   i64.const 1
	   i32.const 1024
	   i64.extend_u/i32
	   tee_local $l0
	   i64.const 4
	   call $ask_external_data)
	 (func $execute (export "execute") (type $t0)
	   (local $idx i32)
	   (set_local $idx (i32.const 0))
	   (block
		(loop
		(set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
		(br_if 0 (i32.lt_u (get_local $idx) (i32.const %d)))
		)
	   )
	   i32.const 1024
	   i64.extend_u/i32
	   i64.const 4
	   call $set_return_data)
	 (table $T0 1 1 anyfunc)
	 (memory $memory (export "memory") 17)
	 (data (i32.const 1024) "new beeb"))`, time, time)))
}
