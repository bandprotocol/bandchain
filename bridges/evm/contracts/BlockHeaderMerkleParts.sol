pragma solidity 0.5.14;
import {Utils} from "./Utils.sol";

/// @dev Library for computing Tendermint's block header hash from app hash, time, and height.
///
/// In Tendermint, a block header hash is the Merkle hash of a binary tree with 16 leaf nodes.
/// Each node encodes a data piece of the blockchain. The notable data leaves are: [C] app_hash,
/// [2] height, and [3] - time. All data pieces are combined into one 32-byte hash to be signed
/// by block validators. The structure of the Merkle tree is shown below.
///
///                                   [BlockHeader]
///                                /                \
///                   [3A]                                    [3B]
///                 /      \                                /      \
///         [2A]                [2B]                [2C]                [2D]
///        /    \              /    \              /    \              /    \
///    [1A]      [1B]      [1C]      [1D]      [1E]      [1F]      [1G]      [1H]
///    /  \      /  \      /  \      /  \      /  \      /  \      /  \      /  \
///  [0]  [1]  [2]  [3]  [4]  [5]  [6]  [7]  [8]  [9]  [A]  [B]  [C]  [D]  [E]  [F]
///
///  [0] - version   [1] - chain_id          [2] - height                [3] - time
///  [4] - num_txs   [5] - total_txs         [6] - last_block_id         [7] - last_commit_hash
///  [8] - data_hash [9] - validators_hash   [A] - next_validators_hash  [B] - consensus_hash
///  [C] - app_hash  [D] - last_results_hash [E] - evidence_hash         [F] - proposer_address
///
/// Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
/// root hash, since we only want to validate the correctness of [C] and [2]. In fact, only
/// [1A], [3], [2B], [2C], [D], and [1H] are needed in order to compute [BlockHeader].
library BlockHeaderMerkleParts {
    struct Data {
        bytes32 versionAndChainIdHash; // [1A]
        bytes32 timeHash; // [3]
        bytes32 txCountAndLastBlockInfoHash; // [2B]
        bytes32 consensusDataHash; // [2C]
        bytes32 lastResultsHash; // [D]
        bytes32 evidenceAndProposerHash; // [1H]
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
                            Utils.merkleLeafHash(
                                Utils.encodeVarintUnsigned(_blockHeight)
                            ), // [2]
                            _self.timeHash
                        )
                    ), // [3]
                    _self.txCountAndLastBlockInfoHash
                ), // [2B]
                Utils.merkleInnerHash( // [3B]
                    _self.consensusDataHash, // [2C]
                    Utils.merkleInnerHash( // [2D]
                        Utils.merkleInnerHash( // [1G]
                            Utils.merkleLeafHash(
                                abi.encodePacked(uint8(32), _appHash)
                            ), // [C]
                            _self.lastResultsHash
                        ), // [D]
                        _self.evidenceAndProposerHash
                    )
                )
            ); // [1H]
    }
}
