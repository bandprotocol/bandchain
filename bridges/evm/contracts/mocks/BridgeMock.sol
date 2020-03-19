pragma solidity 0.5.14;
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
    Bridge.VerifyOracleDataResult public latestResult;
    IBridge public bridge;

    constructor(IBridge _bridge) public {
        bridge = _bridge;
    }

    function relayAndSafe(bytes calldata _data) external {
        latestResult = bridge.relayAndVerify(_data);
    }
}
