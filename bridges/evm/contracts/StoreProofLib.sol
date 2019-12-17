pragma solidity 0.5.14;

import { BytesLib } from "./BytesLib.sol";

/**
 * @dev Helper library for calculating Merkle proof of Tendermint's multistore and iAVL tree.
 */
library StoreProofLib {
  using BytesLib for uint256;
  using BytesLib for bytes32;

  /**
   * @dev A group of data that necessary for computing appHash
   * @param prefixs is an array of unit256 (prefix) which is represented in the following format
   * +---------------------------------------+
   * | L/R | Length |     ...     |   Data   | -> Data is the output of this function
   * +---------------------------------------+
   * |  1  |    7   | 248-Length  | 8*Length | -> Represent size (bit)
   * +---------------------------------------+
   * The first bit (L/R) is used to represent whether the current path is a left child of an AVL tree
   * If first bit (L/R) is 0 then the current path is a right child
   * If first bit (L/R) is 1 then the current path is a left child
   * @param paths is an array of bytes32 (path) that represented an associated hash which is needed to compute root hash
   * @param otherMSHashes (other multistore hashes) is a hashing of 8 difference multistore hashes
   * The hashing schemes is definded as the following ascii
   *                                            ____________appHash____________
   *                                          /                                \
   *                   _________ otherMSHashes _________                        \
   *                 /                                  \                        \
   *         _____ h5 ______                      ______ h6 _______               \
   *       /                \                   /                  \               \
   *     h1                  h2               h3                    h4              \
   *     /\                  /\               /\                    /\               \
   *  acc  distribution   gov  main     params  slashing     staking  supply          zoracle
   *
   * Notice that all mutistore names are sorted lexically
   * @param key is a request storage prefix (1 bytes) + requestID (8 bytes) in zoracle mutistore
   * So the size of key is 9 bytes
   * @param value is a data that stored in Tendermint storage
   */
  struct Data {
    uint256[] prefixes;
    bytes32[] paths;
    bytes32 otherMSHashes;
    uint64 key;
    bytes value;
  }

  /**
   * @dev Returns the hash of a Merkle leaf node in a store
   * This function will be reverted if there is no prefix, because the first prefix is always a prefix of leaf hash
   * The "uint8(1)" is the prefix of request storage within zoracle store which is "0x01"
   * The "uint8(9)" is represented the length of its posfix which in this case is key and its prefix ( sp.key and request storage prefix )
   * The "uint8(32)" is represented the length of its posfix which in this case is "sha256(sp.value)"
   */
  function getLeafHash(Data memory sp) internal pure returns(bytes32) {
    require(sp.prefixes.length > 0, "FIRST_PREFIX_IS_NEEDED");
    return sha256(abi.encodePacked(
      sp.prefixes[0].getBytes(),
      uint8(9),
      uint8(1),
      uint64(sp.key),
      uint8(32),
      sha256(sp.value)
    ));
  }

  /**
   * @dev Returns a computed store hash
   * By following Tendermint hashing schemes, the root hash of AVL tree need double sha256 hashing at the end
   * The prefixes's size need to be paths's size + 1, because first prefix is for leaf hash calculation and others are associate with paths
   */
  function getAVLHash(Data memory sp) internal pure returns(bytes32) {
    require(sp.prefixes.length == sp.paths.length + 1, "LENGTH_OF_PREFIXS_AND_PATHS_ARE_INCOMPATIBLE");
    bytes32 h = getLeafHash(sp);

    for (uint256 i = 1; i < sp.prefixes.length; i++) {
      if (sp.prefixes[i] >> 255 == 1) {
        h = sha256(abi.encodePacked(sp.prefixes[i].getBytes(),uint8(32),sp.paths[i-1],uint8(32),h));
      } else {
        h = sha256(abi.encodePacked(sp.prefixes[i].getBytes(),uint8(32),h,uint8(32),sp.paths[i-1]));
      }
    }

    return sha256(abi.encodePacked(sha256(abi.encodePacked(h))));
  }

  /**
   * @dev Returns an app hash
   */
  function getAppHash(Data memory sp) internal pure returns(bytes32) {
    return sp.otherMSHashes.innerHash(
      sha256(abi.encodePacked(uint8(0), uint8(7), "zoracle", uint8(32), getAVLHash(sp)))
    );
  }
}
