// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;
import {Bridge, IBridge} from "../Bridge.sol";

/// @dev Mock OracleBridge that allows setting oracle iAVL state at a given height directly.
contract BridgeMock is Bridge {
    constructor(ValidatorWithPower[] memory _validators)
        public
        Bridge(_validators)
    {}

    function setOracleState(uint256 _blockHeight, bytes32 _oracleIAVLStateHash)
        public
    {
        oracleStates[_blockHeight] = _oracleIAVLStateHash;
    }
}

contract ReceiverMock {
    Bridge.RequestPacket public latestReq;
    Bridge.ResponsePacket public latestRes;
    Bridge.Packet[] public latestPackets;
    IBridge public bridge;

    constructor(IBridge _bridge) public {
        bridge = _bridge;
    }

    function relayAndSafe(bytes calldata _data) external {
        (latestReq, latestRes) = bridge.relayAndVerify(_data);
    }

    function relayAndMultiSafe(bytes calldata _data) external {
        Bridge.Packet[] memory packets = bridge.relayAndMultiVerify(_data);
        delete latestPackets;
        for (uint256 i = 0; i < packets.length; i++) {
            latestPackets.push(packets[i]);
        }
    }
}
