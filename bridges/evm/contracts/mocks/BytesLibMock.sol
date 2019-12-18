pragma solidity 0.5.14;

import { BytesLib } from "../BytesLib.sol";

contract BytesLibMock{
  function leafHash(bytes memory value) public pure returns (bytes32) {
    return BytesLib.leafHash(value);
  }

  function innerHash(bytes32 left, bytes32 right) public pure returns (bytes32) {
    return BytesLib.innerHash(left, right);
  }

  function decodeVarint(bytes memory encodeByte) public pure returns (uint256) {
    return BytesLib.decodeVarint(encodeByte);
  }

  function getBytes(uint256 _prefix) public pure returns(bytes memory) {
    return BytesLib.getBytes(_prefix);
  }

  function getSegment(bytes memory bs, uint256 start, uint256 end) public pure returns(bytes memory) {
    return BytesLib.getSegment(bs, start, end);
  }
}
