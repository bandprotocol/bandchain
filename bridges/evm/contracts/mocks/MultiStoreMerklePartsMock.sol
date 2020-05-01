pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {MultiStore} from "../MultiStore.sol";


contract MultiStoreMock {
    function getAppHash(MultiStore.Data memory _self)
        public
        pure
        returns (bytes32)
    {
        return MultiStore.getAppHash(_self);
    }
}
