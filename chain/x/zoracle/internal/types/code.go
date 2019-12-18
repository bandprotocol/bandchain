package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// StoredCode store actual code with owner
type StoredCode struct {
	Code  []byte
	Owner sdk.AccAddress
}

// NewStoredCode is a constructor of StoredCode
func NewStoredCode(code []byte, owner sdk.AccAddress) StoredCode {
	return StoredCode{
		Code:  code,
		Owner: owner,
	}
}

// GetCodeHash is a function to calculate hash of actual code
func (c StoredCode) GetCodeHash() []byte {
	return crypto.Sha256(c.Code)
}
