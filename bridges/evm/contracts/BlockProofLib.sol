pragma solidity 0.5.14;

import { BytesLib } from "./BytesLib.sol";

library BlockProofLib {
  using BytesLib for bytes;
  using BytesLib for bytes32;

  struct Data {
    bytes32 appHash;
    bytes encodedHeight;
    bytes32[] others;
    bytes leftMsg;
    uint256 rightMsgSeperator;
    bytes rightMsg;
    bytes signatures;
  }

  struct SignaturesTracker {
    uint256 noSig;
    address lastSigner;
  }

  function calculateBlockHash(Data memory bp) internal pure returns(bytes32) {
    require(bp.others.length == 6, "PROOF_SIZE_MUST_BE_6");
    bytes32 left = bp.encodedHeight.leafHash().innerHash(bp.others[1]);
    left = bp.others[0].innerHash(left);
    left = left.innerHash(bp.others[2]);
    bytes32 right = (abi.encodePacked(uint8(32), bp.appHash)).leafHash().innerHash(bp.others[4]);
    right = right.innerHash(bp.others[5]);
    right = bp.others[3].innerHash(right);
    return left.innerHash(right);
  }

  function getSignBytes(Data memory bp, uint256 start, uint256 end) internal pure returns(bytes32) {
    bytes32 blockHash = calculateBlockHash(bp);
    return sha256(abi.encodePacked(bp.leftMsg, blockHash, bp.rightMsg.getSegment(start,end)));
  }

  function getSignersFromSignatures(Data memory bp) internal pure returns(address[] memory) {
    SignaturesTracker memory st;
    bytes32 r;
    bytes32 s;
    uint8 v;

    st.lastSigner = address(0);
    bytes memory signatures = bp.signatures;
    uint256 seperator = bp.rightMsgSeperator;

    // Verify signature with signBytes
    require(signatures.length % 65 == 0, "INVALID_SIGNATURE_LENGTH");
    st.noSig = signatures.length / 65;

    require(st.noSig == seperator >> 248, "INCOMPATIBLE_SEPERATOR_AND_SIGS");

    address[] memory signers = new address[](st.noSig);
    uint256 accl = 0;
    for (uint256 i = 0; i < st.noSig; i++) {
      assembly {
        r := mload(add(signatures, add(mul(65, i), 32)))
        s := mload(add(signatures, add(mul(65, i), 64)))
        v := and(mload(add(signatures, add(mul(65, i), 65))), 255)
      }
      if (v < 27) {
        v += 27;
      }
      require(v == 27 || v == 28, "INVALID_SIGNATURE");

      signers[i] = ecrecover(getSignBytes(bp, accl, accl + (seperator & 255)), v, r, s);
      require(st.lastSigner < signers[i], "SIG_ORDER_INVALID");

      st.lastSigner = signers[i];
      accl += seperator & 255;
      seperator >>= 8;
    }
    return signers;
  }
}
