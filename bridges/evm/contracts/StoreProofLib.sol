pragma solidity 0.5.14;

import { BytesLib } from "./BytesLib.sol";

/**
 * @dev Helper library for calculating Merkle proof
 */
library StoreProofLib {
  using BytesLib for uint256;

  /**
   * @dev A group of data that necessary for computing appHash
   */
  struct Data {
    uint256[] prefixes;
    bytes32[] path;
    bytes32 otherMSHashes;
    uint64 key;
    bytes value;
  }

  /**
   * @dev Returns the hash of a Merkle leaf node in a store
   */
  function getLeafHash(Data memory sp) internal pure returns(bytes32) {
    require(sp.prefixes.length > 0);
    return sha256(abi.encodePacked(
      sp.prefixes[0].getBytes(),
      uint8(9),
      uint8(1),
      uint64(sp.key),
      uint8(32),
      sha256(abi.encodePacked(sp.value))
    ));
  }

  /**
   * @dev Returns a computed store hash
   */
  function getAVLHash(Data memory sp) internal pure returns(bytes32) {
    require(sp.prefixes.length == sp.path.length + 1);
    bytes32 h = getLeafHash(sp);

    for (uint256 i = 1; i < sp.prefixes.length; i++) {
      if (sp.prefixes[i] >> 255 == 1) {
        h = sha256(abi.encodePacked(sp.prefixes[i].getBytes(),uint8(32),sp.path[i-1],uint8(32),h));
      } else {
        h = sha256(abi.encodePacked(sp.prefixes[i].getBytes(),uint8(32),h,uint8(32),sp.path[i-1]));
      }
    }

    return sha256(abi.encodePacked(sha256(abi.encodePacked(h))));
  }

  /**
   * @dev Returns an app hash
   */
  function getAppHash(Data memory sp) internal pure returns(bytes32) {
    return sha256(abi.encodePacked(
      uint8(1),
      sp.otherMSHashes,
      sha256(abi.encodePacked(uint8(0), uint8(7), "zoracle", uint8(32), getAVLHash(sp)))
    ));
  }
}
