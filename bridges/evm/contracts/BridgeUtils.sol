pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {IBridge} from "./IBridge.sol";


library BridgeUtils {
    /// Returns the hash of a RequestPacket.
    /// @param _request A tuple that represents RequestPacket struct.
    function getRequestKey(IBridge.RequestPacket memory _request)
        internal
        pure
        returns (bytes32)
    {
        return keccak256(abi.encode(_request));
    }
}
