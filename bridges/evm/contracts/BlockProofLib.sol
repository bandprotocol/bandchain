pragma solidity 0.5.14;

import { BytesLib } from "./BytesLib.sol";

library BlockProofLib {
  using BytesLib for bytes;

  struct Data {
    bytes32 appHash;
    bytes encodedHeight;
    bytes32[] others;
    bytes leftMsg;
    bytes rightMsg;
    bytes signatures;
  }

  function calculateBlockHash(Data memory bp) internal pure returns(bytes memory) {
    require(bp.others.length == 6, "PROOF_SIZE_MUST_BE_6");
    bytes memory left = bp.encodedHeight.leafHash().innerHash(abi.encodePacked(bp.others[1]));
    left = (abi.encodePacked(bp.others[0])).innerHash(left);
    left = left.innerHash(abi.encodePacked(bp.others[2]));
    bytes memory right = (abi.encodePacked(hex"20", bp.appHash)).leafHash().innerHash(abi.encodePacked(bp.others[4]));
    right = right.innerHash(abi.encodePacked(bp.others[5]));
    right = (abi.encodePacked(bp.others[3])).innerHash(right);
    return left.innerHash(right);
  }

  function getSignersFromSignatures(Data memory bp) internal pure returns(address[] memory) {
    bytes memory blockHash = calculateBlockHash(bp);
    bytes32 signBytes = sha256(abi.encodePacked(bp.leftMsg, blockHash, bp.rightMsg));
    address lastSigner = address(0);
    bytes32 r;
    bytes32 s;
    uint8 v;
    
    bytes memory signatures = bp.signatures;

    // Verify signature with signBytes
    require(signatures.length % 65 == 0, "INVALID_SIGNATURE_LENGTH");
    uint256 noSig = signatures.length / 65;

    address[] memory signers = new address[](noSig);

    for (uint i = 0; i < noSig; i++) {
      assembly {
        r := mload(add(signatures, add(mul(65, i), 32)))
        s := mload(add(signatures, add(mul(65, i), 64)))
        v := and(mload(add(signatures, add(mul(65, i), 65))), 255)
      }
      if (v < 27) {
        v += 27;
      }
      require(v == 27 || v == 28, "INVALID_SIGNATURE");
      signers[i] = ecrecover(signBytes, v, r, s);
      require(lastSigner < signers[i], "SIG_ORDER_INVALID");
      lastSigner = signers[i];
    }
    return signers;
  }
}
