package proof

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

func TestBlockHeaderMerkleParts(t *testing.T) {
	layout := "2006-01-02T15:04:05.000000000Z"
	str := "2020-04-20T03:30:30.143851745Z"
	blockTime, _ := time.Parse(layout, str)

	// Copy block header Merkle Part here
	header := types.Header{
		Version: version.Consensus{Block: 10, App: 0},
		ChainID: "bandchain",
		Height:  381837,
		Time:    blockTime,
		LastBlockID: types.BlockID{
			Hash: hexToBytes("F633B30D4FBEC862F4A041311E2CB7DFAD63D57930B065A563299449D25BD9CE"),
			PartsHeader: types.PartSetHeader{
				Total: 1,
				Hash:  hexToBytes("7F334B7EE4F8AAC5E70F07FEB9A58A72F120E9AC046167851FC94BC4F2729550"),
			},
		},
		LastCommitHash:     hexToBytes("561D0BB2B6A6E58E20A6BED0F16C8FF5E333BB5A93C69A8E7F3C13542A84DB60"),
		DataHash:           nil,
		ValidatorsHash:     hexToBytes("3AEB137B43144B229F0CA7AC43033E03FCEE25877A3661E88848E436C3D6DD65"),
		NextValidatorsHash: hexToBytes("3AEB137B43144B229F0CA7AC43033E03FCEE25877A3661E88848E436C3D6DD65"),
		ConsensusHash:      hexToBytes("AD82B220C509602720D74FD75BCE7CFE9B148039958F236D8894E00EB1599E04"),
		AppHash:            hexToBytes("1CCD765C80D0DC1705BB7B6BE616DAD3CF2E6439BB9A9B776D5BD183F89CA141"),
		LastResultsHash:    nil,
		EvidenceHash:       nil,
		ProposerAddress:    hexToBytes("F23391B5DBF982E37FB7DADEA64AAE21CAE4C172"),
	}
	blockMerkleParts := GetBlockHeaderMerkleParts(amino.NewCodec(), &header)
	expectBlockHash := hexToBytes("A35617A81409CE46F1F820450B8AD4B217D99AE38AAA719B33C4FC52DCA99B22")
	appHash := hexToBytes("1CCD765C80D0DC1705BB7B6BE616DAD3CF2E6439BB9A9B776D5BD183F89CA141")
	blockHeight := 381837

	// Verify code
	blockHash := branchHash(
		branchHash(
			branchHash(
				blockMerkleParts.VersionAndChainIdHash,
				branchHash(
					leafHash(cdcEncode(amino.NewCodec(), blockHeight)),
					blockMerkleParts.TimeHash,
				),
			),
			blockMerkleParts.LastBlockIDAndOther,
		),
		branchHash(
			branchHash(
				blockMerkleParts.NextValidatorHashAndConsensusHash,
				branchHash(
					leafHash(cdcEncode(amino.NewCodec(), appHash)),
					blockMerkleParts.LastResultsHash,
				),
			),
			blockMerkleParts.EvidenceAndProposerHash,
		),
	)
	require.Equal(t, expectBlockHash, blockHash)
}

func TestBlockHeaderMerklePartsSecond(t *testing.T) {
	layout := "2006-01-02T15:04:05.000000000Z"
	str := "2020-04-22T18:46:58.678092672Z"
	blockTime, _ := time.Parse(layout, str)

	// Copy block header Merkle Part here
	header := types.Header{
		Version: version.Consensus{Block: 10, App: 0},
		ChainID: "bandchain",
		Height:  664,
		Time:    blockTime,
		LastBlockID: types.BlockID{
			Hash: hexToBytes("08CBE9C9127D469983BD12774A36C298EEA960960BDFADDA5988CF8EE869EC59"),
			PartsHeader: types.PartSetHeader{
				Total: 1,
				Hash:  hexToBytes("090942025A2E8C0DAFAF8A5EAAC5837248BD0247A26E012EA9CC0302657C0877"),
			},
		},
		LastCommitHash:     hexToBytes("7F465E59CE6AF9BC6F78C43F510288715CB5DC268E8EFD6DF3357DA27ED9AC8A"),
		DataHash:           nil,
		ValidatorsHash:     hexToBytes("3AEB137B43144B229F0CA7AC43033E03FCEE25877A3661E88848E436C3D6DD65"),
		NextValidatorsHash: hexToBytes("3AEB137B43144B229F0CA7AC43033E03FCEE25877A3661E88848E436C3D6DD65"),
		ConsensusHash:      hexToBytes("AD82B220C509602720D74FD75BCE7CFE9B148039958F236D8894E00EB1599E04"),
		AppHash:            hexToBytes("2F3BEAC1586C205052B74E1CF3D284CD022F739200B74CB51B910D0F3D0BF13D"),
		LastResultsHash:    nil,
		EvidenceHash:       nil,
		ProposerAddress:    hexToBytes("F0C23921727D869745C4F9703CF33996B1D2B715"),
	}
	blockMerkleParts := GetBlockHeaderMerkleParts(amino.NewCodec(), &header)
	expectBlockHash := hexToBytes("DDABA9F8C3BB8B22CC94E1768E56E4E6111DE2BC33ED76D3CC5EB73748C89D3E")
	appHash := hexToBytes("2F3BEAC1586C205052B74E1CF3D284CD022F739200B74CB51B910D0F3D0BF13D")
	blockHeight := 664

	// Verify code
	blockHash := branchHash(
		branchHash(
			branchHash(
				blockMerkleParts.VersionAndChainIdHash,
				branchHash(
					leafHash(cdcEncode(amino.NewCodec(), blockHeight)),
					blockMerkleParts.TimeHash,
				),
			),
			blockMerkleParts.LastBlockIDAndOther,
		),
		branchHash(
			branchHash(
				blockMerkleParts.NextValidatorHashAndConsensusHash,
				branchHash(
					leafHash(cdcEncode(amino.NewCodec(), appHash)),
					blockMerkleParts.LastResultsHash,
				),
			),
			blockMerkleParts.EvidenceAndProposerHash,
		),
	)
	require.Equal(t, expectBlockHash, blockHash)
}
