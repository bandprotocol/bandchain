package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataSource is the data structure for storing data sources in the storage.
type DataSource struct {
	Owner       sdk.AccAddress `json:"owner"`       // Address authorized to edit this data source.
	Name        string         `json:"name"`        // Short string explaining this data source.
	Description string         `josn:"description"` // Longer string explaining what this does.
	Fee         sdk.Coins      `json:"fee"`         // Amount of fee sent to owner for every query.
	Executable  []byte         `json:"executable"`  // Executable script to be run by validators.
}

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
