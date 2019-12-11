pragma solidity 0.5.14;

import { UintLib } from "./UintLib.sol";

library StoreProofLib {
  using UintLib for uint256;

  struct Data {
    uint256[] prefixes;
    bytes32[] path;
    bytes32 otherMSHashes;
    uint64 key;
    bytes value;
  }

  function getLeafHash(Data memory sp) internal pure returns(bytes32) {
    require(sp.prefixes.length > 0);
    return sha256(abi.encodePacked(
      sp.prefixes[0].toU8Arr(),
      uint8(9),
      uint8(1),
      uint64(sp.key),
      uint8(32),
      sha256(abi.encodePacked(sp.value))
    ));
  }

  function getAVLHash(Data memory sp) internal pure returns(bytes32) {
    require(sp.prefixes.length == sp.path.length + 1);
    bytes32 h = getLeafHash(sp);

    for (uint256 i = 1; i < sp.prefixes.length; i++) {
      if (sp.prefixes[i] >> 255 == 1) {
        h = sha256(abi.encodePacked(sp.prefixes[i].toU8Arr(),uint8(32),sp.path[i-1],uint8(32),h));
      } else {
        h = sha256(abi.encodePacked(sp.prefixes[i].toU8Arr(),uint8(32),h,uint8(32),sp.path[i-1]));
      }
    }

    return h;
  }

  function getAppHash(Data memory sp) internal pure returns(bytes32) {
    bytes32 zoracle = sha256(
      abi.encodePacked(
        sha256(
          abi.encodePacked(getAVLHash(sp))
        )
      )
    );

    return sha256(abi.encodePacked(
      uint8(1),
      sp.otherMSHashes,
      sha256(abi.encodePacked(uint8(0), uint8(7), "zoracle", uint8(32), zoracle))
    ));
  }
}
