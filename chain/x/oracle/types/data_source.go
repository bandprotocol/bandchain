package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataSource is a type to store detail of data source.
type DataSource struct {
	Owner       sdk.AccAddress `json:"owner"`
	Name        string         `json:"name"`
	Description string         `josn:"description"`
	Fee         sdk.Coins      `json:"fee"`
	Executable  []byte         `json:"executable"`
}

// NewDataSource creates a new DataSource instance.
func NewDataSource(owner sdk.AccAddress, name string, description string, fee sdk.Coins, executable []byte) DataSource {
	return DataSource{
		Owner:       owner,
		Name:        name,
		Description: description,
		Fee:         fee,
		Executable:  executable,
	}
}
