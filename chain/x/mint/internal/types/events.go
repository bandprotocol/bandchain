package types

import "github.com/cosmos/cosmos-sdk/x/mint"

// Minting module event types
const (
	EventTypeMint = mint.ModuleName

	AttributeKeyBondedRatio      = "bonded_ratio"
	AttributeKeyInflation        = "inflation"
	AttributeKeyAnnualProvisions = "annual_provisions"
)
