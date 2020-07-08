pragma solidity 0.6.0;
import {Utils} from "./Utils.sol";

// Computes Tendermint's application state hash at this given block. AppHash is actually a
// Merkle hash on muliple stores.
//                         ________________[AppHash]_______________
//                        /                                        \
//             _______[I9]______                          ________[I10]________
//            /                  \                       /                     \
//       __[I5]__             __[I6]__              __[I7]__               __[I8]__
//      /         \          /         \           /         \            /         \
//    [I1]       [I2]     [I3]        [I4]       [8]        [9]          [A]        [B]
//   /   \      /   \    /    \      /    \
// [0]   [1]  [2]   [3] [4]   [5]  [6]    [7]
// [0] - acc      [1] - distr   [2] - evidence  [3] - gov
// [4] - main     [5] - mint    [6] - oracle    [7] - params
// [8] - slashing [9] - staking [A] - supply    [D] - upgrade
// Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
// root hash, since we only want to validate the correctness of [6] In fact, only
// [7], [I3], [I5], and [I10] are needed in order to compute [AppHash].

library MultiStore {
    struct Data {
        bytes32 accToGovStoresMerkleHash; // [I5]
        bytes32 mainAndMintStoresMerkleHash; // [I3]
        bytes32 oracleIAVLStateHash; // [6]
        bytes32 paramsStoresMerkleHash; // [7]
        bytes32 slashingToUpgradeStoresMerkleHash; // [I10]
    }

    function getAppHash(Data memory _self) internal pure returns (bytes32) {
        return
            Utils.merkleInnerHash( // [AppHash]
                Utils.merkleInnerHash( // [I9]
                    _self.accToGovStoresMerkleHash, // [I5]
                    Utils.merkleInnerHash( // [I6]
                        _self.mainAndMintStoresMerkleHash, // [I3]
                        Utils.merkleInnerHash(
                            Utils.merkleLeafHash( // [I4]
                                abi.encodePacked( // [6]
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
                            ),
                            _self.paramsStoresMerkleHash // [7]
                        )
                    )
                ),
                _self.slashingToUpgradeStoresMerkleHash // [I10]
            );
    }
}
