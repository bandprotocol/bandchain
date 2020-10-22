package common

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EvMap is a type alias for SDK events mapping from Attr.Key to the list of values.
type EvMap map[string][]string

// JsDict is a type alias for JSON dictionary.
type JsDict map[string]interface{}

// Message is a simple wrapper data type for each message published to Kafka.
type Message struct {
	Key   string
	Value JsDict
}

// atoi converts the given string into an int64. Panics on errors.
func Atoi(val string) int64 {
	res, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}

// atoi converts the given string into an uint64. Panics on errors.
func Atoui(val string) uint64 {
	res, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		panic(err)
	}
	return res
}

func ParseEvents(events sdk.StringEvents) EvMap {
	evMap := make(EvMap)
	for _, event := range events {
		for _, kv := range event.Attributes {
			key := event.Type + "." + kv.Key
			evMap[key] = append(evMap[key], kv.Value)
		}
	}
	return evMap
}
