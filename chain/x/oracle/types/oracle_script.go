package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewOracleScript creates a new OracleScript instance.
func NewOracleScript(
	owner sdk.AccAddress, name string, description string, code []byte, schema string, sourceCodeURL string) OracleScript {
	return OracleScript{
		Owner:         owner,
		Name:          name,
		Description:   description,
		Code:          code,
		Schema:        schema,
		SourceCodeURL: sourceCodeURL,
	}
}
