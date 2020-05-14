package owasm

import (
	"github.com/perlin-network/life/compiler"
)

const UnknownGasCost = 10000000

var defaultGas = map[string]int64{
	// Registers
	"get_local":  3,
	"get_global": 3,
	"set_local":  3,
	"set_global": 3,

	// Memory
	"i32.load8_s":  3,
	"i32.load8_u":  3,
	"i32.load16_s": 3,
	"i32.load16_u": 3,
	"i32.load":     3,
	"i64.load8_s":  3,
	"i64.load8_u":  3,
	"i64.load16_s": 3,
	"i64.load16_u": 3,
	"i64.load32_s": 3,
	"i64.load32_u": 3,
	"i64.load":     3,
	"f32.load":     3,
	"f64.load":     3,
	"memory.size":  100, // Breaking out of the VM
	"memory.grow":  100, // Breaking out of the VM

	"i32.store":   3,
	"i32.store8":  3,
	"i32.store16": 3,
	"i64.store":   3,
	"i64.store8":  3,
	"i64.store16": 3,
	"i64.store32": 3,
	"f32.store":   3,
	"f64.store":   3,

	// Flow  Control
	"nop":         0,
	"unreachable": 0,
	"jmp":         2,
	"jmp_if":      2, // Estimate equal to jmp
	"jmp_table":   2, // It's  the same opcode of br_table
	"return":      2,

	// Calls
	"call":          2,
	"call_indirect": 100, // Breaking out of the VM

	// Constants
	"i32.const": 0,
	"i64.const": 0,
	"f32.const": 0,
	"f64.const": 0,

	// 32-bit Integer operators
	"i32.add":    1,
	"i32.sub":    1,
	"i32.mul":    3,
	"i32.div_s":  80,
	"i32.div_u":  80,
	"i32.rem_s":  80,
	"i32.rem_u":  80,
	"i32.and":    1,
	"i32.or":     1,
	"i32.xor":    1,
	"i32.shl":    1,
	"i32.shr_u":  1,
	"i32.shr_s":  1,
	"i32.rotl":   2,
	"i32.rotr":   2,
	"i32.eq":     1,
	"i32.eqz":    1,
	"i32.ne":     1,
	"i32.lt_s":   1,
	"i32.lt_u":   1,
	"i32.le_s":   1,
	"i32.le_u":   1,
	"i32.gt_s":   1,
	"i32.gt_u":   1,
	"i32.ge_s":   1,
	"i32.ge_u":   1,
	"i32.clz":    105,
	"i32.ctz":    105,
	"i32.popcnt": 4,

	// 64-bit Integer operators
	"i64.add":    1,
	"i64.sub":    1,
	"i64.mul":    3,
	"i64.div_s":  80,
	"i64.div_u":  80,
	"i64.rem_s":  80,
	"i64.rem_u":  80,
	"i64.and":    1,
	"i64.or":     1,
	"i64.xor":    1,
	"i64.shl":    1,
	"i64.shr_u":  1,
	"i64.shr_s":  1,
	"i64.rotl":   2,
	"i64.rotr":   2,
	"i64.eq":     1,
	"i64.eqz":    1,
	"i64.ne":     1,
	"i64.lt_s":   1,
	"i64.lt_u":   1,
	"i64.le_s":   1,
	"i64.le_u":   1,
	"i64.gt_s":   1,
	"i64.gt_u":   1,
	"i64.ge_s":   1,
	"i64.ge_u":   1,
	"i64.clz":    UnknownGasCost,
	"i64.ctz":    UnknownGasCost,
	"i64.popcnt": 4,

	// 32-bit Float operators
	"f32.add":      1,
	"f32.sub":      1,
	"f32.mul":      9,
	"f32.div":      80,
	"f32.min":      1,
	"f32.max":      1,
	"f32.copysign": 1,
	"f32.eq":       1,
	"f32.ne":       1,
	"f32.lt":       1,
	"f32.le":       1,
	"f32.gt":       1,
	"f32.ge":       1,
	"f32.sqrt":     80,
	"f32.ceil":     1,
	"f32.floor":    1,
	"f32.trunc":    1,
	"f32.nearest":  1,
	"f32.abs":      1,
	"f32.neg":      1,

	// 64-bit Float operators
	"f64.add":      2,
	"f64.sub":      2,
	"f64.mul":      9,
	"f64.div":      80,
	"f64.min":      1,
	"f64.max":      1,
	"f64.copysign": 1,
	"f64.eq":       1,
	"f64.ne":       1,
	"f64.lt":       1,
	"f64.le":       1,
	"f64.gt":       1,
	"f64.ge":       1,
	"f64.sqrt":     80,
	"f64.ceil":     1,
	"f64.floor":    1,
	"f64.trunc":    1,
	"f64.nearest":  1,
	"f64.abs":      1,
	"f64.neg":      1,

	// Datatype conversions, truncations, reinterpretations, promotions, and demotions
	"i32.wrap/i64":        3,
	"i64.extend_s/i32":    3, //This cost equals the sum of i64.shl and i64.shr_s
	"i64.extend_u/i32":    3,
	"i32.reinterpret/f32": 1,
	"i64.reinterpret/f64": 1,
	"f32.reinterpret/i32": 1,
	"f64.reinterpret/i64": 1,
	"i32.trunc_u/f32":     3,
	"i32.trunc_u/f64":     3,
	"i64.trunc_u/f32":     3,
	"i64.trunc_u/f64":     3,
	"i32.trunc_s/f32":     3,
	"i32.trunc_s/f64":     3,
	"i64.trunc_s/f32":     3,
	"i64.trunc_s/f64":     3,
	"f32.demote/f64":      3,
	"f64.promote/f32":     3,
	"f32.convert_u/i32":   3,
	"f32.convert_u/i64":   3,
	"f64.convert_u/i32":   3,
	"f64.convert_u/i64":   3,
	"f32.convert_s/i32":   3,
	"f32.convert_s/i64":   3,
	"f64.convert_s/i32":   3,
	"f64.convert_s/i64":   3,

	// Type-parametric operators.
	"drop":   3,
	"select": 3,
}

type BandChainGasPolicy struct{}

func (p *BandChainGasPolicy) GetCost(key compiler.Instr) int64 {
	gasCost, found := defaultGas[key.Op]
	if !found {
		return UnknownGasCost
	}
	return gasCost
}
