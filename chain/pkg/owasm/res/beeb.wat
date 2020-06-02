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


;; (module
;;   (import "env" "requestExternalData" (func (param i64 i64 i32 i64) (result i64)))
;;   (import "env" "saveReturnData" (func (param i32 i64) (result i64)))
;;   (func (;"prepare": Requests external data from sources 1, 2, 3 with call data "beeb";)
;;     i64.const 1 (;data source id;)
;;     i64.const 1 (;external id;)
;;     i32.const 1048576 (;a raw pointer of string;)
;;     i64.const 4 (;string length;)
;;     call 0 (;call function requestExternalData;)
;;     drop (;clear stack;)
;;     i64.const 2 (;data source id;)
;;     i64.const 2 (;external id;)
;;     i32.const 1048576 (;a raw pointer of string;)
;;     i64.const 4 (;string length;)
;;     call 0 (;call function requestExternalData;)
;;     drop (;clear stack;)
;;     i64.const 3 (;data source id;)
;;     i64.const 3 (;external id;)
;;     i32.const 1048576 (;a raw pointer of string;)
;;     i64.const 4  (;string length;)
;;     call 0 (;call function requestExternalData;)
;;     drop (;clear stack;)
;;   )
;;   (func (;"execute": Resolves with result "beeb";)
;;     i32.const 1048576 (;a raw pointer of string;)
;;     i64.const 4 (;string length;)
;;     call 1 (;call function saveReturnData;)
;;     drop (;clear stack;)
;;   )
;;   (memory 17)
;;   (data (i32.const 1048576) "beeb") (;str = "beeb";)
;;   (export "prepare" (func 2))
;;   (export "execute" (func 3)))
