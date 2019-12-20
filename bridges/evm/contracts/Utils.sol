pragma solidity 0.5.14;


/// @dev Helper utility library for calculating Merkle proof and managing bytes.
library Utils {
  /// @dev Returns the hash of a Merkle leaf node.
  function merkleLeafHash(bytes memory _value) internal pure returns (bytes32) {
    return sha256(abi.encodePacked(uint8(0), _value));
  }

  /// @dev Returns the hash of internal node, calculated from child nodes.
  function merkleInnerHash(bytes32 _left, bytes32 _right) internal pure returns (bytes32) {
    return sha256(abi.encodePacked(uint8(1), _left, _right));
  }

  /// @dev Returns the encoded bytes using signed varint encoding of the given input.
  function encodeVarintSigned(uint256 _value) internal pure returns (bytes memory) {
    return encodeVarintUnsigned(_value*2);
  }

  /// @dev Returns the encoded bytes using unsigned varint encoding of the given input.
  function encodeVarintUnsigned(uint256 _value) internal pure returns (bytes memory) {
    // Computes the size of the encoded value.
    uint256 tempValue = _value;
    uint256 size = 0;
    while (tempValue > 0) {
      ++size;
      tempValue >>= 7;
    }
    // Allocates the memory buffer and fills in the encoded value.
    bytes memory result = new bytes(size);
    tempValue = _value;
    for (uint256 idx = 0; idx < size; ++idx) {
      result[idx] = byte(uint8(128) | uint8(tempValue & 127));
      tempValue >>= 7;
    }
    result[size-1] &= byte(uint8(127));  // Drop the first bit of the last byte.
    return result;
  }
}
