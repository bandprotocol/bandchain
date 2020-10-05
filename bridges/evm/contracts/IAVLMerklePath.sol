// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.6.11;
import {Utils} from "./Utils.sol";

/// @dev Library for computing iAVL Merkle root from (1) data leaf and (2) a list of "MerklePath"
/// from such leaf to the root of the tree. Each Merkle path (i.e. proof component) consists of:
///
/// - isDataOnRight: whether the data is on the right subtree of this internal node.
/// - subtreeHeight: well, it is the height of this subtree.
/// - subtreeVersion: the latest block height that this subtree has been updated.
/// - siblingHash: 32-byte hash of the other child subtree
///
/// To construct a hash of an internal Merkle node, the hashes of the two subtrees are combined
/// with extra data of this internal node. See implementation below. Repeatedly doing this from
/// the leaf node until you get to the root node to get the final iAVL Merkle hash.
library IAVLMerklePath {
    struct Data {
        bool isDataOnRight;
        uint8 subtreeHeight;
        uint256 subtreeSize;
        uint256 subtreeVersion;
        bytes32 siblingHash;
    }

    /// @dev Returns the upper Merkle hash given a proof component and hash of data subtree.
    /// @param _dataSubtreeHash The hash of data subtree up until this point.
    function getParentHash(Data memory _self, bytes32 _dataSubtreeHash)
        internal
        pure
        returns (bytes32)
    {
        bytes32 leftSubtree = _self.isDataOnRight
            ? _self.siblingHash
            : _dataSubtreeHash;
        bytes32 rightSubtree = _self.isDataOnRight
            ? _dataSubtreeHash
            : _self.siblingHash;
        return
            sha256(
                abi.encodePacked(
                    _self.subtreeHeight << 1, // Tendermint signed-int8 encoding requires multiplying by 2
                    Utils.encodeVarintSigned(_self.subtreeSize),
                    Utils.encodeVarintSigned(_self.subtreeVersion),
                    uint8(32), // Size of left subtree hash
                    leftSubtree,
                    uint8(32), // Size of right subtree hash
                    rightSubtree
                )
            );
    }
}
