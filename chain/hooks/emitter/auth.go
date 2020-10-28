package emitter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"

	"github.com/bandprotocol/bandchain/chain/hooks/common"
)

func (h *Hook) emitAuthModule(ctx sdk.Context) {
	h.accountKeeper.IterateAccounts(ctx, func(account exported.Account) bool {
		h.Write("SET_ACCOUNT", common.JsDict{
			"address": account.GetAddress(),
			"balance": h.bankKeeper.GetCoins(ctx, account.GetAddress()).String(),
		})
		return false
	})
}
