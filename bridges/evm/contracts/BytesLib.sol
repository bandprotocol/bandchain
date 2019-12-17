pragma solidity 0.5.14;

/**
 * @dev Helper library for calculating Merkle proof and managing bytes.
 */
library BytesLib {
  /**
   * @dev Returns the hash of a Merkle leaf node
   */
  function leafHash(bytes memory value) internal pure returns (bytes32) {
    return sha256(abi.encodePacked(uint8(0), value));
  }

  /**
   * @dev Returns the hash of internal node, calculated from child nodes
   */
  function innerHash(bytes32 left, bytes32 right) internal pure returns (bytes32) {
    return sha256(abi.encodePacked(uint8(1), left, right));
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
   * @dev Returns data part from input prefix (from tendermint), prefix is represented in the following format
   * +---------------------------------------+
   * | L/R | Length |     ...     |   Data   | -> Data is the output of this function
   * +---------------------------------------+
   * |  1  |    7   | 248-Length  | 8*Length | -> Represent size (bit)
   * +---------------------------------------+
   * The first bit not relevant here (used to specify whether a node is left or right child on other contexts).
   * The next 7 bits encode the size of the data in bytes, using big-endian.
   * The last size bytes of the given input are the actual data.
   */
  function getBytes(uint _prefix) internal pure returns (bytes memory) {
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
