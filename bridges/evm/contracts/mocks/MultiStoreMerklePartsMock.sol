pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {MultiStoreMerkleParts} from "../MultiStoreMerkleParts.sol";


contract MultiStoreMerklePartsMock {
    function getAppHash(
        MultiStoreMerkleParts.Data memory _self,
        bytes memory _oraclePrefix
    ) public pure returns (bytes32) {
        return MultiStoreMerkleParts.getAppHash(_self, _oraclePrefix);
    }
}
