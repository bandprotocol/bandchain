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

    function requestOracle(IBridge.RequestPacket calldata request)
        external
        returns (bytes32)
    {
        emit OracleRequest(
            request.clientId,
            request.oracleScriptId,
            request.params,
            request.askCount,
            request.minCount
        );
        bytes32 key = keccak256(abi.encodePacked(msg.sender, request.clientId));
        require(!pendingRequest[key].isValue, "DUPLICATE_PENDING_REQUEST");
        IBridge.RequestPacket memory temp = request;
        pendingRequest[key] = PendingRequest(temp, msg.sender, true);
        return key;
    }
}
