pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;
import { Bridge } from "../Bridge.sol";

/// @dev Mock OracleBridge that allows setting oracle iAVL state at a given height directly.
contract BridgeMock is Bridge {
  constructor(address[] memory _validators) public Bridge(_validators) {}
  function setOracleState(uint256 _blockHeight, bytes32 _oracleIAVLStateHash) public {
    oracleStates[_blockHeight] = _oracleIAVLStateHash;
  }
}
