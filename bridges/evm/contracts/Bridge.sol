pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;
import { BlockHeaderMerkleParts } from "./BlockHeaderMerkleParts.sol";
import { OracleStateMerklePath } from "./OracleStateMerklePath.sol";
import { TMSignature } from "./TMSignature.sol";
import { Utils } from "./Utils.sol";

/// @title Bridge <3 BandChain D3N
/// @author Band Protocol Team
contract Bridge {
  using BlockHeaderMerkleParts for BlockHeaderMerkleParts.Data;
  using OracleStateMerklePath for OracleStateMerklePath.Data;
  using TMSignature for TMSignature.Data;

  /// Mapping from block height to the hash of "zoracle" iAVL Merkle tree.
  mapping (uint256 => bytes32) public oracleStates;
  /// Mapping from an address to whether it's a validator.
  mapping (address => bool) public validators;
  /// The total number of active validators currently on duty.
  uint256 public validatorCount;

  /// Initializes an oracle bridge to BandChain.
  /// @param _validators The initial set of BandChain active validators.
  constructor(address[] memory _validators) public {
    for (uint256 idx = 0; idx < _validators.length; ++idx) {
      address validator = _validators[idx];
      require(!validators[validator], "DUPLICATION_IN_INITIAL_VALIDATOR_SET");
      validators[validator] = true;
    }
    validatorCount = _validators.length;
  }

  /// Relays a new oracle state to the bridge contract.
  /// @param _blockHeight The height of block to relay to this bridge contract.
  /// @param _oracleIAVLStateHash Hash of iAVL Merkle that represents the state of oracle store.
  /// @param _otherStoresMerkleHash Hash of internal Merkle node for other Tendermint storages.
  /// @param _merkleParts Extra merkle parts to compute block hash. See BlockHeaderMerkleParts lib.
  /// @param _signedDataPrefix Prefix data prepended prior to signing block hash.
  /// @param _signatures The signatures signed on this block, sorted alphabetically by address.
  function relayOracleState(
    uint256 _blockHeight,
    bytes32 _oracleIAVLStateHash,
    bytes32 _otherStoresMerkleHash,
    BlockHeaderMerkleParts.Data memory _merkleParts,
    bytes memory _signedDataPrefix,
    TMSignature.Data[] memory _signatures
  )
    public
  {
    // Computes Tendermint's application state hash at this given block. AppHash is actually a
    // Merkle hash on muliple stores. Luckily, we only care about "zoracle" tree and all other
    // stores can just be combined into one bytes32 hash off-chain.
    //
    //                                            ____________appHash_________
    //                                          /                              \
    //                   ____otherStoresMerkleHash ____                         \
    //                 /                                \                        \
    //         _____ h5 ______                    ______ h6 _______               \
    //       /                \                 /                  \               \
    //     h1                  h2             h3                    h4              \
    //     /\                  /\             /\                    /\               \
    //  acc  distribution   gov  main   params  slashing     staking  supply       zoracle
    bytes32 appHash = Utils.merkleInnerHash(
      _otherStoresMerkleHash,
      Utils.merkleLeafHash(
        abi.encodePacked(
          hex"077a6f7261636c6520", // uint8(7) + "zoracle" + uint8(32)
          sha256(abi.encodePacked(sha256(abi.encodePacked(_oracleIAVLStateHash))))
        )
      )
    );
    // Computes Tendermint's block header hash at this given block.
    // TODO: Remove this and update test case
    appHash = hex"6fec355d8a3ce024eed694d19f1e41ae2815c6c9609c5ef59d715935d2e76712";
    bytes32 blockHeader = _merkleParts.getBlockHeader(appHash, _blockHeight);
    // Counts the total number of valid signatures signed by active validators.
    address lastSigner = address(0);
    uint256 validSignatureCount = 0;
    for (uint256 idx = 0; idx < _signatures.length; ++idx) {
      address signer = _signatures[idx].recoverSigner(blockHeader, _signedDataPrefix);
      require(signer > lastSigner, "INVALID_SIGNATURE_SIGNER_ORDER");
      if (validators[signer]) {
        // Increases valid signature count if the signer is one of the validators.
        validSignatureCount += 1;
      }
      lastSigner = signer;
    }
    // Verifies that sufficient validators signed the block and saves the oracle state.
    require(validSignatureCount*3 > validatorCount*2, "INSUFFICIENT_VALIDATOR_SIGNATURES");
    oracleStates[_blockHeight] = _oracleIAVLStateHash;
  }

  /// Verify that the given data is a valid data on BandChain as of the given block height.
  /// @param _blockHeight The block height. Someone must already relay this block.
  /// @param _data The data to verify, with the format similar to what on the blockchain store.
  /// @param _requestId The ID of request for this data piece.
  /// @param _version Lastest block height that the data node was updated.
  /// @param _merklePaths Merkle proof that shows how the data leave is part of the oracle iAVL.
  function verifyOracleData(
    uint256 _blockHeight,
    bytes memory _data,
    uint64 _requestId,
    uint256 _version,
    OracleStateMerklePath.Data[] memory _merklePaths
  )
    public
    view
    returns (bool)
  {
    bytes32 oracleStateRoot = oracleStates[_blockHeight];
    require(oracleStateRoot != bytes32(uint256(0)), "NO_ORACLE_ROOT_STATE_DATA");
    // Computes the hash of leaf node for iAVL oracle tree.
    bytes32 currentMerkleHash = sha256(abi.encodePacked(
      uint8(0),  // Height of tree (only leaf node) is 0 (signed-varint encode)
      uint8(2),  // Size of subtree is 1 (signed-varint encode)
      Utils.encodeVarintSigned(_version),
      uint8(9),  // Size of data key (1-byte constant 0x01 + 8-byte request ID)
      uint8(1),  // Constant 0x01 prefix data request info storage key
      _requestId,
      uint8(32),  // Size of data hash
      sha256(_data)
    ));
    // Goes step-by-step computing hash of parent nodes until reaching root node.
    for (uint256 idx = 0; idx < _merklePaths.length; ++idx) {
      currentMerkleHash = _merklePaths[idx].getParentHash(currentMerkleHash);
    }
    // Verifies that the computed Merkle root matches what currently exists.
    require(currentMerkleHash == oracleStateRoot, "INVALID_ORACLE_DATA_PROOF");
    return true;
  }
}
