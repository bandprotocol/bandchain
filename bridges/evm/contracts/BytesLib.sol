pragma solidity 0.5.14;

/**
 * @dev Helper library using in calculating merkle proof
 */
library BytesLib {
  /**
    * @dev Returns the hash of leaf node of merkle tree
  */
  function leafHash(bytes memory value) internal pure returns (bytes memory) {
    return abi.encodePacked(sha256(abi.encodePacked(false, value)));
  }

  /**
    * @dev Returns the hash of internal node calculate from children nodes
  */
  function innerHash(bytes memory left, bytes memory right) internal pure returns (bytes memory) {
    return abi.encodePacked(sha256(abi.encodePacked(true, left, right)));
  }

  /**
    * @dev Returns the decoded uint from input bytes without checking valid varint format
  */
  function decodeVarint(bytes memory _encodeByte) internal pure returns (uint) {
    uint v = 0;
    for (uint i = 0; i < _encodeByte.length; i++) {
      v = v | uint((uint8(_encodeByte[i]) & 127)) << (i*7);
    }
    return v;
  }

  /**
    * @dev Returns data part from input prefix (from tendermint), prefix is represeneted in this form
    * +--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    * | L/R | Length |      ...      |   Data   | -> Data is an output of this function
    * +--+--+--+--+--+--+--+--+--+--+--+--+--+--+
    * |  1  |    7   |  248-Length  |   Length  | -> Represent size (bit)
    * +--+--+--+--+--+--+--+--+--+--+--+--+--+--+
  */
  function getBytes(uint _prefix) internal pure returns(bytes memory) {
    uint prefix = _prefix;
    uint length = (prefix >> 248) & 127;
    bytes memory arr = new bytes(length);
    for (uint i = length; i > 0; i--) {
      arr[i - 1] = byte(uint8(prefix & 255));
      prefix >>= 8;
    }
    return arr;
  }
}
