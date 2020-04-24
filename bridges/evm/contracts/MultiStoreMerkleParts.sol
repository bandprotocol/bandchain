pragma solidity 0.5.14;
import {Utils} from "./Utils.sol";


library MultiStoreMerkleParts {
    struct Data {
        bytes32 accToMainStoresMerkleHash;
        bytes32 mintStoresMerkleHash;
        bytes oracleIAVLStateHash;
        bytes32 paramsAndSlashingStoresMerkleHash;
        bytes32 stakingToUpgradeStoresMerkleHash;
    }

    function getAppHash(Data memory _self, bytes memory _oraclePrefix)
        internal
        pure
        returns (bytes32)
    {
        return
            Utils.merkleInnerHash(
                _self.accToMainStoresMerkleHash,
                Utils.merkleInnerHash(
                    Utils.merkleInnerHash(
                        Utils.merkleInnerHash(
                            _self.mintStoresMerkleHash,
                            Utils.merkleLeafHash(
                                abi.encodePacked(
                                    _oraclePrefix,
                                    sha256(
                                        abi.encodePacked(
                                            sha256(_self.oracleIAVLStateHash)
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
