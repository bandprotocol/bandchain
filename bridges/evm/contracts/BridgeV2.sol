pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {Bridge} from "./Bridge.sol";

import {IBridgeV2} from "./IBridgeV2.sol";
import {IBridge} from "./IBridge.sol";


contract BridgeV2 is Bridge, IBridgeV2 {
    struct PendingRequest {
        IBridge.RequestPacket request;
        address targetAddress;
        bool isValue;
    }

    mapping(bytes32 => PendingRequest) public pendingRequest;
    event RequestOracle(
        string clientId,
        uint256 oracleScriptId,
        string params,
        uint64 askCount,
        uint64 minCount
    );

    function requestOracle(IBridge.RequestPacket calldata request)
        external
        returns (bytes32)
    {
        emit RequestOracle(
            request.clientId,
            request.oracleScriptId,
            request.params,
            request.askCount,
            request.minCount
        );
        bytes32 key = keccak256(abi.encodePacked(msg.sender, request.clientId));
        if (pendingRequest[key].isValue) revert();
        IBridge.RequestPacket memory r = IBridge.RequestPacket(
            request.clientId,
            request.oracleScriptId,
            request.params,
            request.askCount,
            request.minCount
        );
        pendingRequest[key] = PendingRequest(r, msg.sender, true);
        return key;
    }
}
