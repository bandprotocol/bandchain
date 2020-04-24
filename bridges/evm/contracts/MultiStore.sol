pragma solidity 0.5.14;
import {Utils} from "./Utils.sol";


// Computes Tendermint's application state hash at this given block. AppHash is actually a
// Merkle hash on muliple stores. Luckily, we only care about "zoracle" tree and all other
// stores can just be combined into one bytes32 hash off-chain.
//
//                                            ____________appHash_________
//                                          /                              \
//                   ____otherStoresMerkleHash ____                         ___innerHash___
//                 /                                \                     /                  \
//         _____ h5 ______                    ______ h6 _______        supply              zoracle
//       /                \                 /                  \
//     h1                  h2             h3                    h4
//     /\                  /\             /\                    /\
//  acc  distribution   gov  main     mint  params     slashing   staking
library MultiStore {
    struct Data {
        bytes32 accToMainStoresMerkleHash;
        bytes32 mintStoresMerkleHash;
        bytes32 oracleIAVLStateHash;
        bytes32 paramsAndSlashingStoresMerkleHash;
        bytes32 stakingToUpgradeStoresMerkleHash;
    }

    function getAppHash(Data memory _self) internal pure returns (bytes32) {
        return
            Utils.merkleInnerHash(
                _self.accToMainStoresMerkleHash,
                Utils.merkleInnerHash(
                    Utils.merkleInnerHash(
                        Utils.merkleInnerHash(
                            _self.mintStoresMerkleHash,
                            Utils.merkleLeafHash(
                                abi.encodePacked(
                                    hex"066f7261636c6520",
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
                        _self.paramsAndSlashingStoresMerkleHash
                    ),
                    _self.stakingToUpgradeStoresMerkleHash
                )
            );
    }
}
