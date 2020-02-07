package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// StoredCode store actual code with owner
type StoredCode struct {
	Code  cmn.HexBytes   `json:"code"`
	Name  string         `json:"name"`
	Owner sdk.AccAddress `json:"owner"`
}

// NewStoredCode is a constructor of StoredCode
func NewStoredCode(code []byte, name string, owner sdk.AccAddress) StoredCode {
	return StoredCode{
		Code:  code,
		Name:  name,
		Owner: owner,
	}
}

// GetCodeHash is a function to calculate hash of actual code
func (c StoredCode) GetCodeHash() []byte {
	// TODO: Scheme to calculate codehash
	return crypto.Sha256(append([]byte(c.Name), c.Code...))
}
