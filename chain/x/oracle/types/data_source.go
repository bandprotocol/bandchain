package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewDataSource creates a new DataSource instance.
func NewDataSource(
	owner sdk.AccAddress, name string, description string,
	fee sdk.Coins, executable []byte,
) DataSource {
	return DataSource{
		Owner:       owner,
		Name:        name,
		Description: description,
		Fee:         fee,
		Executable:  executable,
	}
}
