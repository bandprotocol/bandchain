pragma solidity 0.5.14;

import { BytesLib } from "../BytesLib.sol";

contract BytesLibMock{
  function leafHash(bytes memory value) public pure returns (bytes memory) {
    return BytesLib.leafHash(value);
  }

  function innerHash(bytes memory left, bytes memory right) public pure returns (bytes memory) {
    return BytesLib.innerHash(left, right);
  }

  function decodeVarint(bytes memory encodeByte) public pure returns (uint) {
    return BytesLib.decodeVarint(encodeByte);
  }

  function getBytes(uint _prefix) public pure returns(bytes memory) {
    return BytesLib.getBytes(_prefix);
  }
}
