pragma solidity 0.5.17;
pragma experimental ABIEncoderV2;

import {Bridge} from "./Bridge.sol";
import {IBridgeV2} from "./BridgeV2.sol";


contract BridgeV2 is Bridge, IBridgeV2 {
    struct Request {
        string clientId;
        uint64 oracleScriptId;
        bytes params;
        int64 askCount;
        int64 minCount;
    }
    struct PendingRequest {
        Request request;
        address targetAddress;
    }

    mapping(bytes32 => PendingRequest) public pendingRequest;
    event RequestOracle(
        string clientId,
        uint256 oracleScriptId,
        bytes params,
        int64 askCount,
        int64 minCount
    );

    function requestOracle(Request calldata request)
        external
        returns (bytes32)
    {
        emit RequestOracle(
            request.clientId,
            request.oracleScriptId,
            request.calldata,
            request.askCount,
            request.minCount
        );
        bytes32 key = keccak256(abi.encodePacked(msg.sender, request.clientId));
        Request memory r = Request(
            request.clientId,
            request.oracleScriptId,
            request.params,
            request.askCount,
            request.minCount
        );
        pendingRequest[key] = PendingRequest(r, msg.sender);
        return key;
    }
}
