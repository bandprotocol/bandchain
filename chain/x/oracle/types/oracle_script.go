package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// OracleScript is a type to store detail of oracle script.
type OracleScript struct {
	Owner         sdk.AccAddress `json:"owner"`       // Address authorized to edit this script.
	Name          string         `json:"name"`        // Short string explaining this oracle script.
	Description   string         `json:"description"` // Longer string explaining what this does.
	Code          []byte         `json:"code"`        // Owasm bytecode to be run on-chain.
	Schema        string         `json:"schema"`
	SourceCodeURL string         `json:"source_code_url"`
}

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
