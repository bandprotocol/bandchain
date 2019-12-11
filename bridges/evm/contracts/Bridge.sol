pragma solidity 0.5.14;
import { ValidatorsManager } from "./ValidatorsManager.sol";

contract Bridge is ValidatorsManager {
  mapping(uint256 => bytes32) public appHashes;

  uint256 public numberOfValidators;
  mapping(address => bool) validators;

  struct localVar {
    bytes blockHash;
    bytes32 signBytes;
    address signer;
    address lastSigner;
    uint noSig;
  }

  constructor(address[] memory _initValidators) public ValidatorsManager(_initValidators) {}

  function toU8Arr(uint256 prefix) public pure returns(bytes memory){
    uint256 n = (prefix >> 248) & 127;
    bytes memory arr = new bytes(n);
    while (n > 0) {
      arr[n-1] = byte(uint8(prefix & 255));
      prefix >>= 8;
      n--;
    }
    return arr;
  }

  function getLeafHash(uint256 firstPrefix, uint64 key, bytes32 valueHash) public pure returns(bytes32) {
    return sha256(abi.encodePacked(
      toU8Arr(firstPrefix),
      uint8(9),
      uint8(1),
      uint64(key),
      uint8(32),
      valueHash
    ));
  }

  function getAVLHash(
    uint256[] memory prefixes,
    bytes32[] memory path,
    uint64 key,
    bytes32 valueHash
  ) public pure returns(bytes32) {
    require(prefixes.length == path.length + 1);
    bytes32 h = getLeafHash(prefixes[0],key,valueHash);

    for (uint256 i = 1; i < prefixes.length; i++) {
      if (prefixes[i] >> 255 == 1) {
        h = sha256(abi.encodePacked(toU8Arr(prefixes[i]),uint8(32),path[i-1],uint8(32),h));
      } else {
        h = sha256(abi.encodePacked(toU8Arr(prefixes[i]),uint8(32),h,uint8(32),path[i-1]));
      }
    }

    return h;
  }

  function getAppHash(
    uint256[] memory prefixes,
    bytes32[] memory path,
    bytes32 otherMSHashes,
    uint64 key,
    bytes memory value
  ) public pure returns(bytes32) {
    bytes32 zoracle = sha256(
      abi.encodePacked(
        sha256(
          abi.encodePacked(getAVLHash(prefixes, path, key, sha256(abi.encodePacked(value))))
        )
      )
    );

    return sha256(abi.encodePacked(
      uint8(1),
      otherMSHashes,
      sha256(abi.encodePacked(uint8(0), uint8(7), "zoracle", uint8(32), zoracle))
    ));
  }

  function verifyAppHash(
    uint256[] memory prefixes,
    bytes32[] memory path,
    bytes32 otherMSHashes,
    uint64 blockHeight,
    uint64 key,
    bytes memory value
  ) public view returns(bool) {
    return appHashes[blockHeight] == getAppHash(prefixes, path, otherMSHashes, key, value);
  }

  function leafHash(bytes memory value) internal pure returns (bytes memory) {
    return abi.encodePacked(sha256(abi.encodePacked(false, value)));
  }

  function innerHash(bytes memory left, bytes memory right) internal pure returns (bytes memory) {
    return abi.encodePacked(sha256(abi.encodePacked(true, left, right)));
  }

  function decodeVarint(bytes memory _encodeByte) public pure returns (uint) {
    uint v = 0;
    for (uint i = 0; i < _encodeByte.length; i++) {
      v = v | uint((uint8(_encodeByte[i]) & 127)) << (i*7);
    }
    return v;
  }

  function calculateBlockHash(
    bytes memory _encodedHeight,
    bytes32 _appHash,
    bytes32[] memory _others
  ) public pure returns(bytes memory) {
    require(_others.length == 6, "PROOF_SIZE_MUST_BE_6");
    bytes memory left = innerHash(leafHash(_encodedHeight), abi.encodePacked(_others[1]));
    left = innerHash(abi.encodePacked(_others[0]), left);
    left = innerHash(left, abi.encodePacked(_others[2]));
    bytes memory right = innerHash(leafHash(abi.encodePacked(hex"20", _appHash)), abi.encodePacked(_others[4]));
    right = innerHash(right, abi.encodePacked(_others[5]));
    right = innerHash(abi.encodePacked(_others[3]), right);
    return innerHash(left, right);
  }

  function verifyValidatorSignatures(
    bytes32 _appHash,
    bytes memory _encodedHeight,
    bytes32[] memory _others,
    bytes memory _leftMsg,
    bytes memory _rightMsg,
    bytes memory _signatures
  ) public view returns(bool) {
    localVar memory vars;
    vars.blockHash = calculateBlockHash(_encodedHeight, _appHash, _others);
    vars.signBytes = sha256(abi.encodePacked(_leftMsg, vars.blockHash, _rightMsg));
    vars.lastSigner = address(0);
    bytes32 r;
    bytes32 s;
    uint8 v;

    // Verify signature with signBytes
    require(_signatures.length % 65 == 0, "INVALID_SIGNATURE_LENGTH");
    vars.noSig = _signatures.length / 65;

    // number of signatures > 2/3 of numberOfValidators
    require(vars.noSig * 3 > numberOfValidators * 2);
    for (uint i = 0; i < vars.noSig; i++) {
      assembly {
        r := mload(add(_signatures, add(mul(65, i), 32)))
        s := mload(add(_signatures, add(mul(65, i), 64)))
        v := and(mload(add(_signatures, add(mul(65, i), 65))), 255)
      }
      if (v < 27) {
        v += 27;
      }
      require(v == 27 || v == 28, "INVALID_SIGNATURE");
      vars.signer = ecrecover(vars.signBytes, v, r, s);
      require(vars.lastSigner < vars.signer, "SIG_ORDER_INVALID");
      require(validators[vars.signer], "INVALID_VALIDATOR_ADDRESS");
      vars.lastSigner = vars.signer;
    }
    return true;
  }

  function submitAppHash(
    bytes32 _appHash,
    bytes memory _encodedHeight,
    bytes32[] memory _others,
    bytes memory _leftMsg,
    bytes memory _rightMsg,
    bytes memory _signatures
  ) public {
    require(verifyValidatorSignatures(_appHash,_encodedHeight,_others,_leftMsg,_rightMsg,_signatures));
    appHashes[decodeVarint(_encodedHeight)] = _appHash;
  }

  function submitAndVerify (
    bytes32 _appHash,
    bytes memory _encodedHeight,
    bytes32[] memory _others,
    bytes memory _leftMsg,
    bytes memory _rightMsg,
    bytes memory _signatures,
    uint256[] memory prefixes,
    bytes32[] memory path,
    bytes32 otherMSHashes,
    uint64 key,
    bytes memory value
  ) public returns(bool) {
    submitAppHash(
      _appHash,
      _encodedHeight,
      _others,
      _leftMsg,
      _rightMsg,
      _signatures
    );
    require(
      verifyAppHash(
        prefixes,
        path,
        otherMSHashes,
        uint64(decodeVarint(_encodedHeight)),
        key,
        value
      ),
      "FAIL_TO_VERIFY_APP_HASH"
    );

    return true;
  }
}
