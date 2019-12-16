pragma solidity 0.5.14;

import { StoreProofLib } from "../StoreProofLib.sol";

contract StoreProofLibMock{
  using StoreProofLib for StoreProofLib.Data;

  function getLeafHash(
    uint256[] memory prefixes,
    bytes32[] memory path,
    bytes32 otherMSHashes,
    uint64 key,
    bytes memory value
  ) public pure returns(bytes32) {
    StoreProofLib.Data memory sp;
    sp.prefixes = prefixes;
    sp.path = path;
    sp.otherMSHashes = otherMSHashes;
    sp.key = key;
    sp.value = value;
    return sp.getLeafHash();
  }

  function getAVLHash(
    uint256[] memory prefixes,
    bytes32[] memory path,
    bytes32 otherMSHashes,
    uint64 key,
    bytes memory value
  ) public pure returns(bytes32) {
    StoreProofLib.Data memory sp;
    sp.prefixes = prefixes;
    sp.path = path;
    sp.otherMSHashes = otherMSHashes;
    sp.key = key;
    sp.value = value;
    return sp.getAVLHash();
  }

  function getAppHash(
    uint256[] memory prefixes,
    bytes32[] memory path,
    bytes32 otherMSHashes,
    uint64 key,
    bytes memory value
  ) public pure returns(bytes32) {
    StoreProofLib.Data memory sp;
    sp.prefixes = prefixes;
    sp.path = path;
    sp.otherMSHashes = otherMSHashes;
    sp.key = key;
    sp.value = value;
    return sp.getAppHash();
  }
}
