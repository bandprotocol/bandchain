package price

import "github.com/bandprotocol/bandchain/chain/x/oracle/types"

type Input struct {
	Symbols    []string `json:"symbols"`
	Multiplier uint64   `json:"multiplier"`
}

type Output struct {
	Pxs []uint64 `json:"pxs"`
}

type Price struct {
	Symbol      string          `json:"symbol"`
	Multiplier  uint64          `json:"multiplier"`
	Px          uint64          `json:"px"`
	RequestID   types.RequestID `json:"request_id"`
	ResolveTime int64           `json:"resolve_time"`
}

func NewPrice(symbol string, multiplier uint64, px uint64, reqID types.RequestID, resolveTime int64) Price {
	return Price{
		Symbol:      symbol,
		Multiplier:  multiplier,
		Px:          px,
		RequestID:   reqID,
		ResolveTime: resolveTime,
	}
}
