// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {BlockHeaderMerkleParts} from "../BlockHeaderMerkleParts.sol";


contract BlockHeaderMerklePartsMock {
    function getBlockHeader(
        BlockHeaderMerkleParts.Data memory _self,
        bytes32 _appHash,
        uint256 _blockHeight
    ) public pure returns (bytes32) {
        return
            BlockHeaderMerkleParts.getBlockHeader(
                _self,
                _appHash,
                _blockHeight
            );
    }
}
