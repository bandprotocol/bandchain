// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Packets} from "./Packets.sol";
import {Bridge} from "./Bridge.sol";
import {ICacheBridge} from "./ICacheBridge.sol";

/// @title CacheBridge <3 BandChain
/// @author Band Protocol Team
contract CacheBridge is Bridge, ICacheBridge {
    /// Mapping from hash of RequestPacket to the latest ResponsePacket.
    mapping(bytes32 => ResponsePacket) public requestsCache;

    /// Initializes an oracle bridge to BandChain by pass the argument to the parent contract (Bridge.sol).
    /// @param _validators The initial set of BandChain active validators.
    constructor(ValidatorWithPower[] memory _validators)
        public
        Bridge(_validators)
    {}

    /// Returns the hash of a RequestPacket.
    /// @param _request A tuple that represents RequestPacket struct.
    function getRequestKey(RequestPacket memory _request)
        public
        view
        override
        returns (bytes32)
    {
        return keccak256(abi.encode(_request));
    }

    /// Returns the ResponsePacket for a given RequestPacket.
    /// Reverts if can't find the related response in the mapping.
    /// @param _request A tuple that represents RequestPacket struct.
    function getLatestResponse(RequestPacket memory _request)
        public
        view
        override
        returns (ResponsePacket memory)
    {
        ResponsePacket memory res = requestsCache[getRequestKey(_request)];
        require(res.requestId != 0, "RESPONSE_NOT_FOUND");

        return res;
    }

    /// Performs oracle state relay and oracle data verification in one go.
    /// After that, the results will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relay(bytes calldata _data) external override {
        (RequestPacket memory req, ResponsePacket memory res) = this
            .relayAndVerify(_data);

        bytes32 requestKey = getRequestKey(req);

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
}
