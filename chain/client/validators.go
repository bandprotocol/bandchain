package rpc

import (
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	secptm "github.com/tendermint/tendermint/crypto/secp256k1"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

type ValidatorMinimal struct {
	Address     string `json:"address"`
	VotingPower int64  `json:"voting_power"`
}

type ValidatorsOnEthereum struct {
	BlockHeight int64              `json:"block_height"`
	Validators  []ValidatorMinimal `json:"validators"`
}

func GetValidators(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		node, err := cliCtx.GetNode()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		validators, err := node.Validators(nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		validatorsOnETH := ValidatorsOnEthereum{}
		validatorsOnETH.BlockHeight = validators.BlockHeight
		validatorsOnETH.Validators = []ValidatorMinimal{}

		for _, validator := range validators.Validators {
			pubKeyBytes, ok := validator.PubKey.(secptm.PubKeySecp256k1)
			if !ok {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, "fail to cast pubkey")
				return
			}

			if pubkey, err := crypto.DecompressPubkey(pubKeyBytes[:]); err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			} else {
				validatorsOnETH.Validators = append(
					validatorsOnETH.Validators,
					ValidatorMinimal{
						Address:     fmt.Sprintf("0x%x", crypto.PubkeyToAddress(*pubkey)),
						VotingPower: validator.VotingPower,
					},
				)
			}
		}

		rest.PostProcessResponseBare(w, cliCtx, validatorsOnETH)
	}
}
