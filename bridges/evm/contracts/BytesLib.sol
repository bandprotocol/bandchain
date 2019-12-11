pragma solidity 0.5.14;

library BytesLib {
  function leafHash(bytes memory value) internal pure returns (bytes memory) {
    return abi.encodePacked(sha256(abi.encodePacked(false, value)));
  }

  function innerHash(bytes memory left, bytes memory right) internal pure returns (bytes memory) {
    return abi.encodePacked(sha256(abi.encodePacked(true, left, right)));
  }

  function decodeVarint(bytes memory _encodeByte) internal pure returns (uint) {
    uint v = 0;
    for (uint i = 0; i < _encodeByte.length; i++) {
      v = v | uint((uint8(_encodeByte[i]) & 127)) << (i*7);
    }
    return v;
  }
}
