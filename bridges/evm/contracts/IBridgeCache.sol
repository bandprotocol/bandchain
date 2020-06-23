pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {IBridge} from "./Bridge.sol";

interface IBridgeCache {
    /// ...
    /// @param _request The ...
    function getLatestResponse(IBridge.RequestPacket calldata _request)
        external
        view
        returns (IBridge.ResponsePacket memory);
}
