pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {Packets} from "./Packets.sol";
import {Bridge} from "./Bridge.sol";

/// @title BridgeWithCache <3 BandChain
/// @author Band Protocol Team
contract BridgeWithCache is Bridge {
    mapping(bytes32 => ResponsePacket) public requestsCache;

    constructor(ValidatorWithPower[] memory _validators)
        public
        Bridge(_validators)
    {}

    function getRequestKey(RequestPacket memory _request)
        public
        pure
        returns (bytes32)
    {
        return keccak256(abi.encode(_request));
    }

    function getLatestResponse(RequestPacket memory _request)
        public
        view
        returns (ResponsePacket memory)
    {
        ResponsePacket memory res = requestsCache[keccak256(
            abi.encode(_request)
        )];
        require(res.requestId != 0, "RESPONSE_NOT_FOUND");

        return res;
    }

    function relay(bytes calldata _data) external {
        (RequestPacket memory req, ResponsePacket memory res) = relayAndVerify(
            _data
        );

        requestsCache[keccak256(abi.encode(req))] = res;
    }
}
