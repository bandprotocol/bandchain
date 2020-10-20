package price

import (
	"github.com/bandprotocol/bandchain/chain/x/oracle/types"
)

type Price struct {
	Symbol      string          `db:"symbol, primarykey" json:"symbol"`
	MinCount    uint64          `db:"min_count, primarykey" json:"min_count"`
	AskCount    uint64          `db:"ask_count, primarykey" json:"ask_count"`
	Multiplier  uint64          `db:"multiplier" json:"multiplier"`
	Px          uint64          `db:"px" json:"px"`
	RequestID   types.RequestID `db:"request_id" json:"request_id"`
	ResolveTime int64           `db:"resolve_id" json:"resolve_time"`
}

func (h *PriceHook) UpsertPrice(price Price) {
	err := h.dbMap.Insert(&price)
	if err != nil {
		h.UpdatePrice(price)
	}
}

func (h *PriceHook) UpdatePrice(price Price) {
	_, err := h.dbMap.Update(&price)
	if err != nil {
		panic(err)
	}
}
