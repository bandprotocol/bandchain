package bridge_test

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	rootmulti "github.com/cosmos/cosmos-sdk/store/rootmulti"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/merkle"
	tmmerkle "github.com/tendermint/tendermint/crypto/merkle"
	tmtypes "github.com/tendermint/tendermint/types"
)

func TestVerify(t *testing.T) {

	multiStoreKey, err := base64.StdEncoding.DecodeString("b3JhY2xl")
	multiStoreData, err := base64.StdEncoding.DecodeString("sQQKrgQKNAoIc2xhc2hpbmcSKAomCIGwChIgg+D8AwR1oaxS4Od9f3ZH1XG8m+1gRaxYJGe3HJX6jJAKMAoEbWFpbhIoCiYIgbAKEiAjmXfb2YDKqg/X5EeV386YpNoE9Zr01JqJ9noPdoLJuwoyCgZvcmFjbGUSKAomCIGwChIgHk8J1UShxanOZPZUuKkKF5pcwUIurM7KaowmfWIoR3wKLwoDZ292EigKJgiBsAoSIGDNE3wZYuysYWOJ1oA0gz8pIVCcLShapfUVOZfOlop0CjIKBnBhcmFtcxIoCiYIgbAKEiC3FYwM8YX5JLzL2tSCZ+JKcbKRnduKueSxdIktWwqCfwovCgNhY2MSKAomCIGwChIg2yfvEwM5jeOsOVuhu7GJjG0s0U3vifgN6gR0AATWPtAKOAoMZGlzdHJpYnV0aW9uEigKJgiBsAoSIBcgbbDh8vDUcxu7sOeBuE97xNZksojsaz9mi6CA4TKTChIKCGV2aWRlbmNlEgYKBAiBsAoKEQoHdXBncmFkZRIGCgQIgbAKCjIKBnN1cHBseRIoCiYIgbAKEiB5vOlNfaFwJr9dyspWHXTvUu0u9P2oO4dRwOTsrvqx1wowCgRtaW50EigKJgiBsAoSIBBkXH0x54SErnFsawHA1BeSX2BHKFxiaY/jys8hTaZrCjMKB3N0YWtpbmcSKAomCIGwChIgAWQd1+QThGQIgplzNoKWduHHnL+lCmSEp0MMFu4i3hI=")
	// multiStoreValue, err := base64.StdEncoding.DecodeString("AAAADGJhbmRjaGFpbi5qcwAAAAAAAAABAAAAEAAAAARCQU5EAAAAAAAPQkAAAAAAAAAABAAAAAAAAAACAAAADGJhbmRjaGFpbi5qcwAAAAAAAAACAAAAAAAAAAQAAAAAXvMr0wAAAABe8yvYAAAAAQAAAAgAAAAAABPQ5g==")
	iavlKey, err := base64.StdEncoding.DecodeString("/wAAAAAAAAAC")
	iavlData, err := base64.StdEncoding.DecodeString("gAYK/QUKLAgkEMawBRiBsAoiIEwLLEW2S35vBmktv/Dl5RpV4fbq2GI89+Z+vpedfbSaCiwIIBDV9wIYgbAKIiB4C1HDQ3u/bQmlKqCT3E+lwq3Yi1MExyCuBavHbTtZdAosCB4Qh7sBGIGwCiIgPM78BG5L8Dkk2wd4GDv62h4btWl7IOD1tGiP2HLq0FEKKwgcEMVzGIGwCiogellLzKQ455Sk1c6F0xw+Qt2Q45uXJySbe9au3Uo7YGsKKwgaEMY+GP6vCiognMTlA8jSJe2gYvHWTSbXJO2xlhxrh4aX+6PmOctZU1QKKwgYEMYeGP6vCiogwFP8iQC+uMZuOrAfCsofJu1eYBjuvuODsdxdk4U/N9cKKwgWEMYOGP6vCiogWKG87K6MxIaXRKm3aYeCoY1ZuvfXP7ycm46czWk/95oKKwgUEOMGGP6vCiogrMJT4wURudOa0O2juC7lBa2itZ630yJvUN6zd2Wd7KwKKwgSEPACGP6vCiogh3NfDJ9Elj7hWT9xQYrUeW/GIycbUZjdoBwZ4ioWxw8KKggQEHoY/q8KIiAdNzzcOVlBusFhtx1J2VtaV9aoryZdwAVKHRYiOmeejQoqCAwQMhj+rwoiIMzEvYwc948XgQoZY7+2YyFs+QOlV1aVmWZz3LdwjjjBCioIChAfGP6vCiogfIpePFoJ5uZiNec7kpAnsRu8MOdMGEtdW8WB8kebZ74KKggIEA8Y/q8KIiBt3MObnVsKDJw7t25na2TZeF+11pAIMQ8JP1wj846V3goqCAYQCBj+rwoiIMCYX8YFsRa5yAkXW115CZ4L23ChWV/U3MIdS4Ws2LVnCioIBBAEGP6vCiIg9SiO/5CeqqCTMfcrlUnOLNHbxcbVZQV3FRlWaW5o9agKKQgCEAIY6xgiIH8HskF3HxoJKvTbymVhpGEwRGP3JmnNnkC/kNw79uSLGjAKCf8AAAAAAAAAAhIgChSHwDyFueESepA4eVyBEPp83N5DiQWuM3S75+kBjc8Y9hM=")

	proof := tmmerkle.Proof{
		Ops: []tmmerkle.ProofOp{
			{
				Type: "iavl:v",
				Key:  iavlKey,
				Data: iavlData,
			},
			{
				Type: "multistore",
				Key:  multiStoreKey,
				Data: multiStoreData,
			},
		},
	}

	prt := rootmulti.DefaultProofRuntime()
	appHash, _ := hex.DecodeString("8607165AC131C8A53EF624E8B86FF8F4E5F27E326FD137F323A0B1A9746C0A21")
	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte("oracle"), merkle.KeyEncodingURL)
	kp = kp.AppendKey(proof.Ops[0].Key, merkle.KeyEncodingURL)
	value, _ := hex.DecodeString("0000000c62616e64636861696e2e6a730000000000000001000000100000000442414e4400000000000f4240000000000000000400000000000000020000000c62616e64636861696e2e6a7300000000000000020000000000000004000000005ef32bd3000000005ef32bd80000000100000008000000000013d0e6")
	err = prt.VerifyValue(&proof, appHash, kp.String(), value)
	require.NoError(t, err)

	cdc := makeCodec()
	proofJson, err := cdc.MarshalJSON(proof)
	require.NoError(t, err)
	fmt.Printf("proofJson: %v\n", string(proofJson))
}

func makeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	authtypes.RegisterCodec(cdc)
	cdc.RegisterConcrete(sdk.TestMsg{}, "cosmos-sdk/Test", nil)
	return cdc
}

func TestRelay(t *testing.T) {
	data := `{ "header": { "version": { "block": "10", "app": "0" }, "chain_id": "band-guanyu-devnet-3", "height": "169986", "time": "2020-06-30T09:52:19.278950714Z", "last_block_id": { "hash": "70DA7DFDB73858905EDEE0119F14882533787D7D986C8E228C7A5D6B6E648029", "parts": { "total": "1", "hash": "CC350C14F6121D10133D6F42DEB0CDCB7BB2136CC8E835873EF624FA594BC461" } }, "last_commit_hash": "D66C13D2211D20AA608FB56990E3A51B1A491B4404155E8CF3A3B25FA5F9FFA6", "data_hash": "", "validators_hash": "BAE324D6C3715BB7E19EBCA90E6DD9195E1E1579C45A2E94B714F6809D052BD5", "next_validators_hash": "BAE324D6C3715BB7E19EBCA90E6DD9195E1E1579C45A2E94B714F6809D052BD5", "consensus_hash": "AD82B220C509602720D74FD75BCE7CFE9B148039958F236D8894E00EB1599E04", "app_hash": "8607165AC131C8A53EF624E8B86FF8F4E5F27E326FD137F323A0B1A9746C0A21", "last_results_hash": "DEB82E155954D6BE14592C66CCF7A1ECE193EEEBCDABAF747B91F44519F09F47", "evidence_hash": "", "proposer_address": "BF2D9CC09716179575BADB6FE65527AA5FE7944C" }, "commit": { "height": "169986", "round": "0", "block_id": { "hash": "A16D5051E98246C13644A0621924F4927D48CF69F798E85ED864159FBE7BAF92", "parts": { "total": "1", "hash": "7C8F4231A587306DBA77BD32E375DBE567CB6B855D0179A4342AF6969C4CC0C9" } }, "signatures": [ { "block_id_flag": 2, "validator_address": "5D1EE10A49A41D89700B3789B75C144B2761494B", "timestamp": "2020-06-30T09:52:22.474499051Z", "signature": "Q3HdhS0ZnIpg9SxybozSeXMST8VtfqERcDjkiTGqtxQvQqwd4NPaC8N+S1MoECVRPLtwKKmfPR/OYgkzuHPzsw==" }, { "block_id_flag": 2, "validator_address": "BF2D9CC09716179575BADB6FE65527AA5FE7944C", "timestamp": "2020-06-30T09:52:22.422707338Z", "signature": "8kNKIWr8lhKZJMGtx7J1G2XRjuqv3oE9L3pq9ny7qO85ejmQlwVrrci7s2ld5sR+LTELfSLp7sj0NsgbmH1q8Q==" }, { "block_id_flag": 2, "validator_address": "BFEA6E94C5DE9D28F38FE44FECAADD6EF1C78683", "timestamp": "2020-06-30T09:52:22.405512992Z", "signature": "m1ul9ZwQIlhr3ZvAFWqnENIT8ym0NtrhVIXKVJMeqCg1jVBtVTWpO5OaVBxmytSanysvpaEPdAdpWIao5PQT0A==" }, { "block_id_flag": 2, "validator_address": "E011D7B2EEB2CFF5119087455E1FC97B97AD5404", "timestamp": "2020-06-30T09:52:22.398992029Z", "signature": "3LOHK86YonGKaQ9/mfmoWjD5IKHpmdQm83fb58rsq/AY1Cq34joI7Zy+8NVGju1ZffQCg20OibuSSRiwWpkrng==" } ] } }`
	var header tmtypes.SignedHeader
	cdc := makeCodec()
	_ = cdc.UnmarshalJSON([]byte(data), &header)
	fmt.Printf("%v\n", header)

	validatorsData := `[
			{
				"address": "5D1EE10A49A41D89700B3789B75C144B2761494B",
				"pub_key": {
					"type": "tendermint/PubKeySecp256k1",
					"value": "AxjvN+RW20yrnfjkb48KuBRu+vDyUnIkKKEX7cLPJ6mo"
				},
				"voting_power": "1001000",
				"proposer_priority": "-77091"
			},
			{
				"address": "BF2D9CC09716179575BADB6FE65527AA5FE7944C",
				"pub_key": {
					"type": "tendermint/PubKeySecp256k1",
					"value": "A3RMvC/B4+gZMBceYrYW/jz1j/FIgymz7xcJMqLO6DYS"
				},
				"voting_power": "1001000",
				"proposer_priority": "954917"
			},
			{
				"address": "BFEA6E94C5DE9D28F38FE44FECAADD6EF1C78683",
				"pub_key": {
					"type": "tendermint/PubKeySecp256k1",
					"value": "Ax0C8EziCltHAJfhkpsQQz/NsrVH+L//tqv4klW7WNv8"
				},
				"voting_power": "1001001",
				"proposer_priority": "1151254"
			},
			{
				"address": "E011D7B2EEB2CFF5119087455E1FC97B97AD5404",
				"pub_key": {
					"type": "tendermint/PubKeySecp256k1",
					"value": "A+/b5PJ6ARfRxG3yQ6fs4Y6TJUjhqrQKM3VS/qqWwKL5"
				},
				"voting_power": "1001000",
				"proposer_priority": "-2029080"
			}
		]`
	var validators []tmtypes.Validator

	_ = cdc.UnmarshalJSON([]byte(validatorsData), &validators)
	fmt.Printf("validatorsData: %v\n", validators)

	totalPowerValidators := int64(0)
	for _, validator := range validators {
		totalPowerValidators += validator.VotingPower
	}

	sumPower := int64(0)

	for idx, signature := range header.Commit.Signatures {
		//TODO: Add logic to check that commit is not use the same validator in multiple time
		for _, validator := range validators {
			if signature.ValidatorAddress.String() == validator.PubKey.Address().String() {
				msg := header.Commit.VoteSignBytes("band-guanyu-devnet-3", idx)
				fmt.Println("msg", msg)
				isValid := validator.PubKey.VerifyBytes(msg, signature.Signature)
				fmt.Println("isValid", isValid)
				if isValid {
					sumPower += validator.VotingPower
				}
				break
			}
		}
	}

	fmt.Println("sumPower", sumPower)
	fmt.Println("totalPowerValidators", totalPowerValidators)

	if 3*sumPower > 2*totalPowerValidators {
		fmt.Println(true)
	} else {
		fmt.Println(false)
	}

	//save validator

	require.Equal(t, 0, 1)
}
