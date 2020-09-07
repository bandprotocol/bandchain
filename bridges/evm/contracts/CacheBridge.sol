// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Packets} from "./Packets.sol";
import {Bridge} from "./Bridge.sol";
import {ICacheBridge} from "./ICacheBridge.sol";
import {BridgeUtils} from "./BridgeUtils.sol";
import {IBridge} from "./IBridge.sol";

/// @title CacheBridge <3 BandChain
/// @author Band Protocol Team
contract CacheBridge is Bridge, ICacheBridge {
    using BridgeUtils for IBridge.RequestPacket;

    /// Mapping from hash of RequestPacket to the latest ResponsePacket.
    mapping(bytes32 => ResponsePacket) public requestsCache;

    /// Initializes an oracle bridge to BandChain by pass the argument to the parent contract (Bridge.sol).
    /// @param _validators The initial set of BandChain active validators.
    constructor(ValidatorWithPower[] memory _validators)
        public
        Bridge(_validators)
    {}

    /// Returns the ResponsePacket for a given RequestPacket.
    /// Reverts if can't find the related response in the mapping.
    /// @param _request A tuple that represents RequestPacket struct.
    function getLatestResponse(RequestPacket memory _request)
        public
        override
        view
        returns (ResponsePacket memory)
    {
        ResponsePacket memory res = requestsCache[_request.getRequestKey()];
        require(res.requestId != 0, "RESPONSE_NOT_FOUND");

        return res;
    }

    /// Performs oracle state relay and oracle data verification in one go.
    /// After that, the results will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relay(bytes calldata _data) external override {
        (RequestPacket memory req, ResponsePacket memory res) = this
            .relayAndVerify(_data);

        bytes32 requestKey = req.getRequestKey();

        require(
            res.resolveTime > requestsCache[requestKey].resolveTime,
            "FAIL_LATEST_REQUEST_SHOULD_BE_NEWEST"
        );

        require(
            res.resolveStatus == 1,
            "FAIL_REQUEST_IS_NOT_SUCCESSFULLY_RESOLVED"
        );

        requestsCache[requestKey] = res;
    }

    /// Performs oracle state relay and many times of oracle data verification in one go.
    /// After that, the results which is an array of Packet will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and an array of data verification.
    function relayMulti(bytes calldata _data) external {
        Packet[] memory packets = this.relayAndMultiVerify(_data);

        for (uint256 i = 0; i < packets.length; i++) {
            Packet memory packet = packets[i];
            bytes32 requestKey = packet.request.getRequestKey();

            require(
                packet.response.resolveTime >
                    requestsCache[requestKey].resolveTime,
                "FAIL_LATEST_REQUEST_SHOULD_BE_NEWEST"
            );

            require(
                packet.response.resolveStatus == 1,
                "FAIL_REQUEST_IS_NOT_SUCCESSFULLY_RESOLVED"
            );

            requestsCache[requestKey] = packet.response;
        }
    }
}
