package keeper

import (
	"github.com/GeoDB-Limited/odincore/chain/x/oracle/types"
)

// modify returns new value if it is not `DoNotModify`. Returns old value otherwise
func modify(oldVal string, newVal string) string {
	if newVal == types.DoNotModify {
		return oldVal
	}
	return newVal
}
