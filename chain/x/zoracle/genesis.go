package zoracle

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	// Scripts []types.StoredCode `json:"scripts"`
}

func NewGenesisState() GenesisState {
	return GenesisState{}
}

func ValidateGenesis(data GenesisState) error {
	// for _, record := range data.RequestRecords {
	// if record.Owner == nil {
	// 	return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Owner", record.Value)
	// }
	// if record.Value == "" {
	// 	return fmt.Errorf("invalid WhoisRecord: Owner: %s. Error: Missing Value", record.Owner)
	// }
	// if record.Price == nil {
	// 	return fmt.Errorf("invalid WhoisRecord: Value: %s. Error: Missing Price", record.Value)
	// }
	// }
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	// for _, storeCode := range data.Scripts {
	// 	codeHash := storeCode.GetCodeHash()
	// 	if keeper.CheckCodeHashExists(ctx, codeHash) {
	// 		continue
	// 	}
	// 	keeper.SetCode(ctx, storeCode.Code, storeCode.Name, storeCode.Owner)
	// }

	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	// var records []Whois
	// iterator := k.GetNamesIterator(ctx)
	// for ; iterator.Valid(); iterator.Next() {

	// 	name := string(iterator.Key())
	// 	whois := k.GetWhois(ctx, name)
	// 	records = append(records, whois)

	// }
	// return GenesisState{WhoisRecords: records}
	return GenesisState{}
}
