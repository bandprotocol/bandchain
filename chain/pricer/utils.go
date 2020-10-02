package pricer

import (
	"strconv"
)

// EvMap is a type alias for SDK events mapping from Attr.Key to the list of values.
type EvMap map[string][]string

type Input struct {
	Symbols    []string `json:"symbols"`
	Multiplier uint64   `json:"multiplier"`
}

type Output struct {
	Pxs []uint64 `json:"pxs"`
}

// atoi converts the given string into an int64. Panics on errors.
func atoi(val string) uint64 {
	res, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}
