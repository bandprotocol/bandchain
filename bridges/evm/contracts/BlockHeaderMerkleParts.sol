pragma solidity 0.6.0;
import {Utils} from "./Utils.sol";


/// @dev Library for computing Tendermint's block header hash from app hash, time, and height.
///
/// In Tendermint, a block header hash is the Merkle hash of a binary tree with 14 leaf nodes.
/// Each node encodes a data piece of the blockchain. The notable data leaves are: [A] app_hash,
/// [2] height, and [3] - time. All data pieces are combined into one 32-byte hash to be signed
/// by block validators. The structure of the Merkle tree is shown below.
///
///                                   [BlockHeader]
///                                /                \
///                   [3A]                                    [3B]
///                 /      \                                /      \
///         [2A]                [2B]                [2C]                [2D]
///        /    \              /    \              /    \              /    \
///    [1A]      [1B]      [1C]      [1D]      [1E]      [1F]        [C]    [D]
///    /  \      /  \      /  \      /  \      /  \      /  \
///  [0]  [1]  [2]  [3]  [4]  [5]  [6]  [7]  [8]  [9]  [A]  [B]
///
///  [0] - version               [1] - chain_id            [2] - height        [3] - time
///  [4] - last_block_id         [5] - last_commit_hash    [6] - data_hash     [7] - validators_hash
///  [8] - next_validators_hash  [9] - consensus_hash      [A] - app_hash      [B] - last_results_hash
///  [C] - evidence_hash         [D] - proposer_address
///
/// Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
/// root hash, since we only want to validate the correctness of [A] and [2]. In fact, only
/// [1A], [3], [2B], [1E], [B], and [2D] are needed in order to compute [BlockHeader].
library BlockHeaderMerkleParts {
    struct Data {
        bytes32 versionAndChainIdHash; // [1A]
        bytes32 timeHash; // [3]
        bytes32 lastBlockIDAndOther; // [2B]
        bytes32 nextValidatorHashAndConsensusHash; // [1E]
        bytes32 lastResultsHash; // [B]
        bytes32 evidenceAndProposerHash; // [2D]
    }

    /// @dev Returns the block header hash after combining merkle parts with necessary data.
    /// @param _appHash The Merkle hash of BandChain application state.
    /// @param _blockHeight The height of this block.
    function getBlockHeader(
        Data memory _self,
        bytes32 _appHash,
        uint256 _blockHeight
    ) internal pure returns (bytes32) {
        return
            Utils.merkleInnerHash( // [BlockHeader]
                Utils.merkleInnerHash( // [3A]
                    Utils.merkleInnerHash( // [2A]
                        _self.versionAndChainIdHash, // [1A]
                        Utils.merkleInnerHash( // [1B]
                            Utils.merkleLeafHash( // [2]
                                Utils.encodeVarintUnsigned(_blockHeight)
                            ),
                            _self.timeHash // [3]
                        )
                    ),
                    _self.lastBlockIDAndOther // [2B]
                ),
                Utils.merkleInnerHash( // [3B]
                    Utils.merkleInnerHash( // [2C]
                        _self.nextValidatorHashAndConsensusHash, // [1E]
                        Utils.merkleInnerHash( // [1F]
                            Utils.merkleLeafHash( // [A]
                                abi.encodePacked(uint8(32), _appHash)
                            ),
                            _self.lastResultsHash // [B]
                        )
                    ),
                    _self.evidenceAndProposerHash // [2D]
                )
            );
    }
}
