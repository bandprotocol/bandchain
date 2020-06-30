package emitter

import (
	"strconv"
)

// EvMap is a type alias for SDK events mapping from Attr.Key to the list of values.
type EvMap map[string][]string

// JsDict is a type alias for JSON dictionary.
type JsDict map[string]interface{}

// atoi converts the given string into an int64. Panics on errors.
func atoi(val string) int64 {
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}
