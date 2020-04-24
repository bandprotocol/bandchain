pragma solidity 0.5.14;
import {Utils} from "./Utils.sol";


// Computes Tendermint's application state hash at this given block. AppHash is actually a
// Merkle hash on muliple stores. Luckily, we only care about "oracle" tree and all other
// stores can just be combined into one bytes32 hash off-chain.
//                         ________________[AppHash]_______________
//                        /                                        \
//             _______[I12]______                         ________[I13]________
//            /                  \                       /                     \
//       __[I8]__             __[I9]__              __[I10]__               __[I11]__
//      /         \          /         \           /         \            /           \
//    [I1]       [I2]     [I3]        [I4]       [I5]        [I6]       [I7]          [E]
//   /   \      /   \    /    \      /    \     /    \       /   \     /    \
// [0]   [1]  [2]   [3] [4]   [5]  [6]    [7] [8]    [9]   [A]   [B] [C]    [D]
// [0] - acc     [1] - bank      [2] - capbility  [3] - distribution  [4] - evidence
// [5] - gov     [6] - ibc       [7] - main       [8] - mint          [9] - oracle
// [A] - params  [B] - slashing  [C] - staking    [D] - supply        [E] - upgrade

library MultiStore {
    struct Data {
        bytes32 accToMainStoresMerkleHash; // [I12]
        bytes32 mintStoresMerkleHash; // [8]
        bytes32 oracleIAVLStateHash; // [9]
        bytes32 paramsAndSlashingStoresMerkleHash; // [I6]
        bytes32 stakingToUpgradeStoresMerkleHash; // [I11]
    }

    function getAppHash(Data memory _self) internal pure returns (bytes32) {
        return
            Utils.merkleInnerHash( // [AppHash]
                _self.accToMainStoresMerkleHash, // [I12]
                Utils.merkleInnerHash( // [I13]
                    Utils.merkleInnerHash( // [I10]
                        Utils.merkleInnerHash( // [I5]
                            _self.mintStoresMerkleHash, // [8]
                            Utils.merkleLeafHash( // [9]
                                abi.encodePacked(
                                    hex"066f7261636c6520", // oracle prefix
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
                    _self.stakingToUpgradeStoresMerkleHash // [I11]
                )
            );
    }
}
