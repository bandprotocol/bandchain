package emitter

import "github.com/cosmos/cosmos-sdk/x/auth/exported"

func (app *App) emitAccountModule() {
	app.AccountKeeper.IterateAccounts(app.DeliverContext, func(account exported.Account) bool {
		app.Write("SET_ACCOUNT", JsDict{
			"address": account.GetAddress(),
			"balance": app.BankKeeper.GetCoins(app.DeliverContext, account.GetAddress()).String(),
		})
		return false
	})
}
