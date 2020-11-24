package proof

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"
)

/*
{
	version: {
		block: "10",
		app: "0"
	},
	chain_id: "band-guanyu-poa",
	height: "3021518",
	time: "2020-11-18T17:43:58.605059026Z",
	last_block_id: {
		hash: "2932C30021FC3AFE1FF0CB7EFDB90BE1CFFA49FAA3B5EFA4FD93DDD54EB2D1E7",
		parts: {
			total: "1",
			hash: "BB5D636F2929F1506301467D612F1D065D3D8F02ABB5E016BE8FEA342B82ACF2"
		}
	},
	last_commit_hash: "95A40FB873116295E3B067F1E8A156F1D4679E538B349DF36F55B96E618FED53",
	data_hash: "F3DC19C19F9620C2E9006C488AAC91F2B5E9A21B5DA8FEA3AB24D4868E558304",
	validators_hash: "4A3CF0246E707DA164809AADCEBCC49B6F238E98FFC7C829004FCADD48B8A63E",
	next_validators_hash: "4A3CF0246E707DA164809AADCEBCC49B6F238E98FFC7C829004FCADD48B8A63E",
	consensus_hash: "0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45",
	app_hash: "91C6C90AD6765C3080CEF2AEB25B1DBDD8ABE6EB409F400C3D6F8DC2767980F6",
	last_results_hash: "6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D",
	evidence_hash: "",
	proposer_address: "A5388FC9C354B02852EE77A87BA0B63155F807EC"
}
*/
func TestBlockHeaderMerkleParts(t *testing.T) {
	// Copy block header Merkle Part here
	header := types.Header{
		Version: version.Consensus{Block: 10, App: 0},
		ChainID: "band-guanyu-poa",
		Height:  3021518,
		Time:    parseTime("2020-11-18T17:43:58.605059026Z"),
		LastBlockID: types.BlockID{
			Hash: hexToBytes("2932C30021FC3AFE1FF0CB7EFDB90BE1CFFA49FAA3B5EFA4FD93DDD54EB2D1E7"),
			PartsHeader: types.PartSetHeader{
				Total: 1,
				Hash:  hexToBytes("BB5D636F2929F1506301467D612F1D065D3D8F02ABB5E016BE8FEA342B82ACF2"),
			},
		},
		LastCommitHash:     hexToBytes("95A40FB873116295E3B067F1E8A156F1D4679E538B349DF36F55B96E618FED53"),
		DataHash:           hexToBytes("F3DC19C19F9620C2E9006C488AAC91F2B5E9A21B5DA8FEA3AB24D4868E558304"),
		ValidatorsHash:     hexToBytes("4A3CF0246E707DA164809AADCEBCC49B6F238E98FFC7C829004FCADD48B8A63E"),
		NextValidatorsHash: hexToBytes("4A3CF0246E707DA164809AADCEBCC49B6F238E98FFC7C829004FCADD48B8A63E"),
		ConsensusHash:      hexToBytes("0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45"),
		AppHash:            hexToBytes("91C6C90AD6765C3080CEF2AEB25B1DBDD8ABE6EB409F400C3D6F8DC2767980F6"),
		LastResultsHash:    hexToBytes("6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D"),
		EvidenceHash:       nil,
		ProposerAddress:    hexToBytes("A5388FC9C354B02852EE77A87BA0B63155F807EC"),
	}
	blockMerkleParts := GetBlockHeaderMerkleParts(amino.NewCodec(), &header)
	expectBlockHash := hexToBytes("707C8E52320D12F0A27A0C35B9DB5649428FC535C4C42A73EFD9BDEB3F8A72D5")
	appHash := hexToBytes("91C6C90AD6765C3080CEF2AEB25B1DBDD8ABE6EB409F400C3D6F8DC2767980F6")

	// Verify code
	blockHash := branchHash(
		branchHash(
			branchHash(
				blockMerkleParts.VersionAndChainIdHash,
				branchHash(
					leafHash(cdcEncode(amino.NewCodec(), header.Height)),
					leafHash(encodeTime(header.Time)),
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
