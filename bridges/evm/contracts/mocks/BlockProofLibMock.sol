pragma solidity 0.5.14;

import { BlockProofLib } from "../BlockProofLib.sol";

library BlockProofLibMock {
  using BlockProofLib for BlockProofLib.Data;

  function calculateBlockHash(
    bytes32 appHash,
    bytes memory encodedHeight,
    bytes32[] memory others,
    bytes memory leftMsg,
    uint256 rightMsgSeperator,
    bytes memory rightMsg,
    bytes memory signatures
  ) public pure returns(bytes32) {
    BlockProofLib.Data memory bp;
    bp.appHash = appHash;
    bp.encodedHeight = encodedHeight;
    bp.others = others;
    bp.leftMsg = leftMsg;
    bp.rightMsgSeperator = rightMsgSeperator;
    bp.rightMsg = rightMsg;
    bp.signatures = signatures;
    return BlockProofLib.calculateBlockHash(bp);
  }

  function getSignersFromSignatures(
    bytes32 appHash,
    bytes memory encodedHeight,
    bytes32[] memory others,
    bytes memory leftMsg,
    uint256 rightMsgSeperator,
    bytes memory rightMsg,
    bytes memory signatures
  ) public pure returns(address[] memory) {
    BlockProofLib.Data memory bp;
    bp.appHash = appHash;
    bp.encodedHeight = encodedHeight;
    bp.others = others;
    bp.leftMsg = leftMsg;
    bp.rightMsgSeperator = rightMsgSeperator;
    bp.rightMsg = rightMsg;
    bp.signatures = signatures;
    return BlockProofLib.getSignersFromSignatures(bp);
  }
}
