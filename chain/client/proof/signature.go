package proof

import (
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	"github.com/tendermint/tendermint/types"
)

type TMSignature struct {
	R                tmbytes.HexBytes `json:"r"`
	S                tmbytes.HexBytes `json:"s"`
	V                uint8            `json:"v"`
	SignedDataSuffix tmbytes.HexBytes `json:"signedDataSuffix"`
}

type TMSignatureEthereum struct {
	R                common.Hash
	S                common.Hash
	V                uint8
	SignedDataSuffix []byte
}

func (signature *TMSignature) encodeToEthFormat() TMSignatureEthereum {
	return TMSignatureEthereum{
		common.BytesToHash(signature.R),
		common.BytesToHash(signature.S),
		signature.V,
		signature.SignedDataSuffix,
	}
}

func recoverETHAddress(msg, sig, signer []byte) ([]byte, uint8, error) {
	for i := uint8(0); i < 2; i++ {
		pubuc, err := crypto.SigToPub(tmhash.Sum(msg), append(sig, byte(i)))
		if err != nil {
			return nil, 0, err
		}
		pub := crypto.CompressPubkey(pubuc)
		var tmp [33]byte

		copy(tmp[:], pub)
		if string(signer) == string(secp256k1.PubKeySecp256k1(tmp).Address()) {
			return crypto.PubkeyToAddress(*pubuc).Bytes(), 27 + i, nil
		}
	}
	return nil, 0, fmt.Errorf("No match address found")
}

func GetSignaturesAndPrefix(info *types.SignedHeader) ([]TMSignature, []byte, error) {
	prefix := ""
	addrs := []string{}
	mapAddrs := map[string]TMSignature{}
	for i, vote := range info.Commit.Signatures {
		if !vote.ForBlock() {
			continue
		}
		msg := info.Commit.VoteSignBytes(info.ChainID, i)
		lr := strings.Split(hex.EncodeToString(msg), hex.EncodeToString(info.Commit.BlockID.Hash))

		if len(lr) != 2 {
			return nil, nil, fmt.Errorf("Split block hash failed")
		}

		if prefix != "" && prefix != lr[0] {
			return nil, nil, fmt.Errorf("Inconsistent prefix signature bytes")
		}
		prefix = lr[0]
		addr, v, err := recoverETHAddress(msg, vote.Signature, vote.ValidatorAddress)
		if err != nil {
			return nil, nil, err
		}
		addrs = append(addrs, string(addr))
		mapAddrs[string(addr)] = TMSignature{
			vote.Signature[:32],
			vote.Signature[32:],
			v,
			mustDecodeString(lr[1]),
		}
	}
	if len(addrs) == 0 {
		return nil, nil, fmt.Errorf("No valid precommit")
	}

	signatures := make([]TMSignature, len(addrs))
	sort.Strings(addrs)
	for i, addr := range addrs {
		signatures[i] = mapAddrs[addr]
	}

	return signatures, mustDecodeString(prefix), nil
}
