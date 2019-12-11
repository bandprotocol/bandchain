package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	trc "github.com/tendermint/tendermint/rpc/core/types"
	trl "github.com/tendermint/tendermint/rpc/lib/types"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/go-amino"
	tc "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

func toByte(hexstr string) []byte {
	if len(hexstr) == 0 {
		return []byte{}
	}
	b, err := hex.DecodeString(hexstr)
	if err != nil {
		panic(err)
	}
	return b
}

func i2b(i int64) []byte {
	s := fmt.Sprintf("%x", i)
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return toByte(s)
}

func GetBlock(blockId uint64) (types.SignedHeader, error) {
	_resp, err := http.Get(fmt.Sprintf("%s/commit?height=%d", strings.Replace(nodeURI, "tcp", "http", 1), blockId))
	if err != nil {
		return types.SignedHeader{}, err
	}
	defer _resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(_resp.Body)
	if err != nil {
		return types.SignedHeader{}, err
	}

	response := trl.RPCResponse{}
	if err := json.Unmarshal(responseBytes, &response); err != nil {
		return types.SignedHeader{}, err
	}

	if response.Error != nil {
		return types.SignedHeader{}, response.Error
	}

	var resp trc.ResultCommit
	if err := (*amino.NewCodec()).UnmarshalJSON(response.Result, &resp); err != nil {
		return types.SignedHeader{}, err
	}

	return resp.SignedHeader, nil
}

func TestBP() {
	block := types.Header{}
	layout := "2006-01-02T15:04:05.000000Z"
	timeStr := "2019-11-28T11:50:14.213447Z"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		panic(err)
	}
	block.Version = version.Consensus{Block: 10, App: 0}
	block.ChainID = "bandchain"
	block.Height = 1900
	block.Time = t
	block.NumTxs = 0
	block.TotalTxs = 0
	block.LastBlockID = types.BlockID{
		Hash: toByte("bad4202a3e08d38d5a0ac08e12d8b88d5e26cf970e61c6a73969291bb0ee8145"),
		PartsHeader: types.PartSetHeader{
			Total: 1,
			Hash:  toByte("3a04d784da77f0406e2dee6b22015a968709320f637893a37484abad51bc7de3"),
		},
	}
	block.LastCommitHash = toByte("d9e35da166023941e8445876e7052ed75214dbf2181413f16d55d97690d602da")
	block.DataHash = []byte{}
	block.ValidatorsHash = toByte("2ce67c22805c7a1f8d9040b76d46038feacf7054a90adba8ecb938c92b88b907")
	block.NextValidatorsHash = toByte("2ce67c22805c7a1f8d9040b76d46038feacf7054a90adba8ecb938c92b88b907")
	block.ConsensusHash = toByte("048091bc7ddc283f77bfbf91d73c44da58c3df8a9cbc867405d8b7f3daada22f")
	block.AppHash = toByte("385bbfe78cc54e09a67caaaf98874f1696f194761361a9f372854ac9d9e23680")
	block.LastResultsHash = []byte{}
	block.EvidenceHash = []byte{}
	block.ProposerAddress = toByte("f23391b5dbf982e37fb7dadea64aae21cae4c172")

	fmt.Println(fmt.Sprintf("appHash -> 0x%x", block.AppHash))

	buf := new(bytes.Buffer)
	amino.EncodeUvarint(buf, uint64(block.Height))
	fmt.Println(fmt.Sprintf("-> 0x%x", buf.Bytes()))

	fmt.Println(hex.EncodeToString(block.Hash()))

	// Run code above get []byte encode by codec and fill merkle
	others := [][]byte{
		merkle.SimpleHashFromByteSlices([][]byte{toByte("080a"), toByte("0962616e64636861696e")}),
		merkle.SimpleHashFromByteSlices([][]byte{toByte("08f6e8feee0510d8e2e365")}),
		merkle.SimpleHashFromByteSlices([][]byte{
			toByte("00"),
			toByte("00"),
			append(
				append(toByte("0a20"), block.LastBlockID.Hash...),
				append(toByte("122408011220"), block.LastBlockID.PartsHeader.Hash...)...,
			),
			append(toByte("20"), block.LastCommitHash...),
		}),
		merkle.SimpleHashFromByteSlices([][]byte{
			[]byte{},
			append(toByte("20"), block.ValidatorsHash...),
			append(toByte("20"), block.NextValidatorsHash...),
			append(toByte("20"), block.ConsensusHash...),
		}),
		merkle.SimpleHashFromByteSlices([][]byte{[]byte{}}),
		merkle.SimpleHashFromByteSlices([][]byte{[]byte{}, append(toByte("14"), block.ProposerAddress...)}),
	}
	s := "["
	for _, h := range others {
		s += fmt.Sprintf(`"0x%s",`, hex.EncodeToString(h))
	}
	s = s[:len(s)-1] + "]"
	fmt.Println(s)
	privS := "a96e62ed3955e65be3aaa3f12d87b6b5cf26039ecfa948dc5107a495418e5430"

	timeStr = "2019-11-28T11:50:19.261931Z"
	t, _ = time.Parse(layout, timeStr)
	vote := types.Vote{
		Type:   2,
		Height: 1900,
		Round:  0,
		BlockID: types.BlockID{
			Hash: block.Hash(),
			PartsHeader: types.PartSetHeader{
				Total: 1,
				Hash:  toByte("8a56cbecfb035c08a4afe57a9a722889c8c6f6a602240ee1803ad4de41936f43"),
			},
		},
		Timestamp: t,
	}
	msg := vote.SignBytes("bandchain")
	lr := strings.Split(hex.EncodeToString(msg), hex.EncodeToString(block.Hash()))
	fmt.Println("0x" + lr[0])
	fmt.Println("0x" + lr[1])
	privB, _ := hex.DecodeString(privS)
	var priv secp256k1.PrivKeySecp256k1
	copy(priv[:], privB)
	pubKey := priv.PubKey()
	pubT, _ := pubKey.(secp256k1.PubKeySecp256k1)
	pub := pubT[:]
	fmt.Println("Pubkey tendermint", hex.EncodeToString(pub))
	fmt.Println("Address tendermint", hex.EncodeToString([]byte(pubKey.Address())))
	sig, err := priv.Sign(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sig1", fmt.Sprintf("%x", sig))
	privE, _ := crypto.HexToECDSA(privS)
	sig, err = crypto.Sign(tc.Sha256(msg), privE)
	// sig, err = crypto.Sign(msg, privE)
	if err != nil {
		panic(err)
	}
	fmt.Println("Pubkey ehtereum", hex.EncodeToString(crypto.FromECDSAPub(&privE.PublicKey)))
	fmt.Println("Ethereum Address", hex.EncodeToString(crypto.PubkeyToAddress(privE.PublicKey).Bytes()))
	fmt.Println("Message hash", hex.EncodeToString(tc.Sha256(msg)))
	// fmt.Println("Message hash", hex.EncodeToString(tc.Sha256(msg)))
	fmt.Println("R:", hex.EncodeToString(sig[:32]))
	fmt.Println("S:", hex.EncodeToString(sig[32:64]))
	fmt.Println("Sig2", fmt.Sprintf("%x", sig))
}

type BlockProof struct {
	EncodedHeight string   `json:"encoded_height"`
	AppHash       string   `json:"app_hash"`
	Others        []string `json:"others"`
	LeftMsg       string   `json:"left_msg"`
	RightMsg      string   `json:"right_msg"`
	Signatures    string   `json:"signatures"`
}

func GetBlockProof(blockId uint64, pk string) (BlockProof, error) {
	bp := BlockProof{}
	sh, err := GetBlock(blockId)
	if err != nil {
		return BlockProof{}, err
	}
	block := *sh.Header
	commit := *sh.Commit

	bp.AppHash = fmt.Sprintf("0x%x", block.AppHash)

	buf := new(bytes.Buffer)
	amino.EncodeUvarint(buf, uint64(block.Height))

	bp.EncodedHeight = fmt.Sprintf("0x%x", buf.Bytes())

	others := [][]byte{
		merkle.SimpleHashFromByteSlices([][]byte{toByte("080a"), toByte("0962616e64636861696e")}),
		merkle.SimpleHashFromByteSlices([][]byte{amino.MustMarshalBinaryBare(block.Time)}),
		merkle.SimpleHashFromByteSlices([][]byte{
			i2b(block.NumTxs),
			i2b(block.TotalTxs),
			append(
				append(toByte("0a20"), block.LastBlockID.Hash...),
				append(toByte("122408011220"), block.LastBlockID.PartsHeader.Hash...)...,
			),
			append(toByte("20"), block.LastCommitHash...),
		}),
		merkle.SimpleHashFromByteSlices([][]byte{
			[]byte{},
			append(toByte("20"), block.ValidatorsHash...),
			append(toByte("20"), block.NextValidatorsHash...),
			append(toByte("20"), block.ConsensusHash...),
		}),
		merkle.SimpleHashFromByteSlices([][]byte{[]byte{}}),
		merkle.SimpleHashFromByteSlices([][]byte{[]byte{}, append(toByte("14"), block.ProposerAddress...)}),
	}

	bp.Others = []string{}
	for _, h := range others {
		bp.Others = append(bp.Others, fmt.Sprintf("0x%s", hex.EncodeToString(h)))
	}

	vote := types.Vote(*commit.Precommits[0])

	msg := vote.SignBytes("bandchain")
	lr := strings.Split(hex.EncodeToString(msg), hex.EncodeToString(block.Hash()))

	bp.LeftMsg = "0x" + lr[0]
	bp.RightMsg = "0x" + lr[1]

	var priv secp256k1.PrivKeySecp256k1
	sig, err := priv.Sign(msg)
	if err != nil {
		return BlockProof{}, err
	}

	privE, _ := crypto.HexToECDSA(pk)
	sig, err = crypto.Sign(tc.Sha256(msg), privE)

	if err != nil {
		return BlockProof{}, err
	}
	bp.Signatures = fmt.Sprintf("0x%x", sig)

	return bp, nil
}
