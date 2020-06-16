(module
  (func (;this function uses 80005 gas for execute;)
    (local $idx i32)
    (set_local $idx (i32.const 0))
    (block
        (loop
          (set_local $idx (get_local $idx) (i32.const 1) (i32.add) )
          (br_if 0 (i32.lt_u (get_local $idx) (i32.const 10000)))
        )
      )
  )
  (export "prepare" (func 0)))
