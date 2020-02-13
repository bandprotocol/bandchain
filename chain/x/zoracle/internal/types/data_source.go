package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataSource is a type to store detail of data source
type DataSource struct {
	Owner      sdk.AccAddress `json:"owner"`
	Name       string         `json:"name"`
	Fee        sdk.Coins      `json:"fee"`
	Executable []byte         `json:"executable"`
}

// NewDataSource - contructor of DataSource struct
func NewDataSource(owner sdk.AccAddress, name string, fee sdk.Coins,	executable []byte) DataSource {
	return DataSource{
		Owner:      owner,
		Name:       name,
		Fee:        fee,
		Executable: executable,
	}
}
