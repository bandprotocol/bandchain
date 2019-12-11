pragma solidity 0.5.14;
import { ValidatorsManager } from "./ValidatorsManager.sol";
import { StoreProofLib } from "./StoreProofLib.sol";
import { BlockProofLib } from "./BlockProofLib.sol";
import { BytesLib } from "./BytesLib.sol";

contract Bridge is ValidatorsManager {
  using BytesLib for bytes;

  using StoreProofLib for StoreProofLib.Data;
  using BlockProofLib for BlockProofLib.Data;

  mapping(uint256 => bytes32) public appHashes;

  constructor(address[] memory _initValidators) public ValidatorsManager(_initValidators) {}

  function verifyAppHash(
    uint64 blockHeight,
    bytes memory value,
    bytes memory storeProof
  ) public view returns(bool) {
    StoreProofLib.Data memory sp;
    sp.value = value;
    (sp.prefixes,sp.path,sp.otherMSHashes,sp.key) = abi.decode(storeProof, (uint256[], bytes32[], bytes32, uint64));

    return appHashes[blockHeight] == sp.getAppHash();
  }

  function submitAppHash(bytes memory appProof) public returns (uint256) {
    BlockProofLib.Data memory bp;
    (bp.appHash,
    bp.encodedHeight,
    bp.others,
    bp.leftMsg,
    bp.rightMsg,
    bp.signatures) = abi.decode(appProof, (bytes32, bytes, bytes32[], bytes, bytes, bytes));

    verifyValidators(bp.getSignersFromSignatures());
    uint256 height = bp.encodedHeight.decodeVarint();
    appHashes[height] = bp.appHash;
    return height;
  }

  function submitAndVerify (
    bytes calldata data, bytes calldata proof
  ) external returns(bool) {
    (bytes memory appProof, bytes memory storeProof) = abi.decode(proof, (bytes, bytes));

    uint256 height = submitAppHash(appProof);
    require(verifyAppHash(uint64(height), data, storeProof), "FAIL_TO_VERIFY_APP_HASH");

    return true;
  }
}
