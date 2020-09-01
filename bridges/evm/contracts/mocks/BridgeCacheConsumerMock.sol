// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;
import {IBridge} from "../IBridge.sol";
import {ICacheBridge} from "../ICacheBridge.sol";

/// @dev Mock Consumer contract for BridgeWithCache which will only consume requests that meet the requirement.
contract BridgeCacheConsumerMock {
    IBridge.RequestPacket public requestTemplate;
    IBridge.ResponsePacket public latestRes;
    ICacheBridge public bridgeCache;

    constructor(
        ICacheBridge _bridgeCache,
        IBridge.RequestPacket memory _initialTemplate
    ) public {
        bridgeCache = _bridgeCache;
        setRequestTemplete(_initialTemplate);
    }

    function setRequestTemplete(IBridge.RequestPacket memory _newTemplate)
        public
    {
        requestTemplate = _newTemplate;
    }

    function consumeCache() public {
        latestRes = bridgeCache.getLatestResponse(requestTemplate);
    }
}
