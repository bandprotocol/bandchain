package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OracleScript is a type to store detail of oracle script.
type OracleScript struct {
	Owner sdk.AccAddress `json:"owner"`
	Name  string         `json:"name"`
	Code  []byte         `json:"code"`
}

// NewOracleScript creates a new OracleScript instance.
func NewOracleScript(owner sdk.AccAddress, name string, code []byte) OracleScript {
	return OracleScript{
		Owner: owner,
		Name:  name,
		Code:  code,
	}
}
