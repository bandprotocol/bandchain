pragma solidity 0.5.14;
import {Utils} from "./Utils.sol";


// Computes Tendermint's application state hash at this given block. AppHash is actually a
// Merkle hash on muliple stores.
//                         ________________[AppHash]_______________
//                        /                                        \
//             _______[I11]______                         ________[I12]________
//            /                  \                       /                     \
//       __[I7]__             __[I8]__              __[I9]__               __[I10]__
//      /         \          /         \           /         \            /           \
//    [I1]       [I2]     [I3]        [I4]       [I5]        [I6]       [C]          [D]
//   /   \      /   \    /    \      /    \     /    \       /   \
// [0]   [1]  [2]   [3] [4]   [5]  [6]    [7] [8]    [9]   [A]   [B]
// [0] - acc     [1] - bank      [2] - capability [3] - distribution  [4] - evidence
// [5] - gov     [6] - ibc       [7] - mem_cap    [8] - mint          [9] - oracle
// [A] - params  [B] - slashing  [C] - staking    [D] - upgrade
// Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
// root hash, since we only want to validate the correctness of [9] In fact, only
// [I11], [8], [9], [I6], and [I10] are needed in order to compute [AppHash].

library MultiStore {
    struct Data {
        bytes32 accToMemCapStoresMerkleHash; // [I11]
        bytes32 mintStoresMerkleHash; // [8]
        bytes32 oracleIAVLStateHash; // [9]
        bytes32 paramsAndSlashingStoresMerkleHash; // [I6]
        bytes32 StakingAndUpgradeStoresMerkleHash; // [I10]
    }

    function getAppHash(Data memory _self) internal pure returns (bytes32) {
        return
            Utils.merkleInnerHash( // [AppHash]
                _self.accToMemCapStoresMerkleHash, // [I11]
                Utils.merkleInnerHash( // [I12]
                    Utils.merkleInnerHash( // [I9]
                        Utils.merkleInnerHash( // [I5]
                            _self.mintStoresMerkleHash, // [8]
                            Utils.merkleLeafHash( // [9]
                                abi.encodePacked(
                                    hex"066f7261636c6520", // oracle prefix (uint8(6) + "oracle" + uint8(32))
                                    sha256(
                                        abi.encodePacked(
                                            sha256(
                                                abi.encodePacked(
                                                    _self.oracleIAVLStateHash
                                                )
                                            )
                                        )
                                    )
                                )
                            )
                        ),
                        _self.paramsAndSlashingStoresMerkleHash // [I6]
                    ),
                    _self.StakingAndUpgradeStoresMerkleHash // [I10]
                )
            );
    }
}
