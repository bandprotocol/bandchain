pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;


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


/// @dev Library for computing Tendermint's block header hash from app hash, time, and height.
///
/// In Tendermint, a block header hash is the Merkle hash of a binary tree with 16 leaf nodes.
/// Each node encodes a data piece of the blockchain. The notable data leaves are: [C] app_hash,
/// [2] height, and [3] - time. All data pieces are combined into one 32-byte hash to be signed
/// by block validators. The structure of the Merkle tree is shown below.
///
///                                   [BlockHeader]
///                                /                \
///                   [3A]                                    [3B]
///                 /      \                                /      \
///         [2A]                [2B]                [2C]                [2D]
///        /    \              /    \              /    \              /    \
///    [1A]      [1B]      [1C]      [1D]      [1E]      [1F]      [1G]      [1H]
///    /  \      /  \      /  \      /  \      /  \      /  \      /  \      /  \
///  [0]  [1]  [2]  [3]  [4]  [5]  [6]  [7]  [8]  [9]  [A]  [B]  [C]  [D]  [E]  [F]
///
///  [0] - version   [1] - chain_id          [2] - height                [3] - time
///  [4] - num_txs   [5] - total_txs         [6] - last_block_id         [7] - last_commit_hash
///  [8] - data_hash [9] - validators_hash   [A] - next_validators_hash  [B] - consensus_hash
///  [C] - app_hash  [D] - last_results_hash [E] - evidence_hash         [F] - proposer_address
///
/// Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
/// root hash, since we only want to validate the correctness of [C] and [2]. In fact, only
/// [1A], [3], [2B], [2C], [D], and [1H] are needed in order to compute [BlockHeader].
library BlockHeaderMerkleParts {
  struct Data {
    bytes32 versionAndChainIdHash;        // [1A]
    bytes32 timeHash;                     // [3]
    bytes32 txCountAndLastBlockInfoHash;  // [2B]
    bytes32 consensusDataHash;            // [2C]
    bytes32 lastResultsHash;              // [D]
    bytes32 evidenceAndProposerHash;      // [1H]
  }

  /// @dev Returns the block header hash after combining merkle parts with necessary data.
  /// @param _appHash The Merkle hash of BandChain application state.
  /// @param _blockHeight The height of this block.
  function getBlockHeader(
    Data memory _self,
    bytes32 _appHash,
    uint256 _blockHeight
  )
    internal
    pure
    returns (bytes32)
  {
    return Utils.merkleInnerHash(                                             // [BlockHeader]
      Utils.merkleInnerHash(                                                  // [3A]
        Utils.merkleInnerHash(                                                // [2A]
          _self.versionAndChainIdHash,                                        // [1A]
          Utils.merkleInnerHash(                                              // [1B]
            Utils.merkleLeafHash(Utils.encodeVarintUnsigned(_blockHeight)),   // [2]
            _self.timeHash)),                                                 // [3]
        _self.txCountAndLastBlockInfoHash),                                   // [2B]
      Utils.merkleInnerHash(                                                  // [3B]
        _self.consensusDataHash,                                              // [2C]
        Utils.merkleInnerHash(                                                // [2D]
          Utils.merkleInnerHash(                                              // [1G]
            Utils.merkleLeafHash(abi.encodePacked(uint8(32), _appHash)),      // [C]
            _self.lastResultsHash),                                           // [D]
          _self.evidenceAndProposerHash)));                                   // [1H]
  }
}

/// @dev Library for performing signer recovery for ECDSA secp256k1 signature. Note that the
/// library is written specifically for signature signed on Tendermint's precommit data, which
/// includes the block hash and some additional information prepended and appended to the block
/// hash. The prepended part (prefix) is the same for all the signers, while the appended part
/// (suffix) is different for each signer (including machine clock, validator index, etc).
library TMSignature {
  struct Data {
    bytes32 r;
    bytes32 s;
    uint8 v;
    bytes signedDataSuffix;
  }

  /// @dev Returns the address that signed on the given block hash.
  /// @param _blockHash The block hash that the validator signed data on.
  /// @param _signedDataPrefix The prefix prepended to block hash before signing.
  function recoverSigner(
    Data memory _self,
    bytes32 _blockHash,
    bytes memory _signedDataPrefix
  )
    internal
    pure
    returns (address)
  {
    return ecrecover(sha256(abi.encodePacked(
      _signedDataPrefix,
      _blockHash,
      _self.signedDataSuffix
    )), _self.v, _self.r, _self.s);
  }
}


/// @dev Library for computing iAVL Merkle root from (1) data leaf and (2) a list of "MerklePath"
/// from such leaf to the root of the tree. Each Merkle path (i.e. proof component) consists of:
///
/// - isDataOnRight: whether the data is on the right subtree of this internal node.
/// - subtreeHeight: well, it is the height of this subtree.
/// - subtreeVersion: the latest block height that this subtree has been updated.
/// - siblingHash: 32-byte hash of the other child subtree
///
/// To construct a hash of an internal Merkle node, the hashes of the two subtrees are combined
/// with extra data of this internal node. See implementation below. Repeatedly doing this from
/// the leaf node until you get to the root node to get the final iAVL Merkle hash.
library OracleStateMerklePath {
  struct Data {
    bool isDataOnRight;
    uint256 subtreeHeight;
    uint256 subtreeSize;
    uint256 subtreeVersion;
    bytes32 siblingHash;
  }

  /// @dev Returns the upper Merkle hash given a proof component and hash of data subtree.
  /// @param _dataSubtreeHash The hash of data subtree up until this point.
  function getParentHash(Data memory _self, bytes32 _dataSubtreeHash)
    internal
    pure
    returns (bytes32)
  {
    bytes32 leftSubtree = _self.isDataOnRight ? _self.siblingHash : _dataSubtreeHash;
    bytes32 rightSubtree = _self.isDataOnRight ? _dataSubtreeHash : _self.siblingHash;
    return sha256(abi.encodePacked(
      Utils.encodeVarintSigned(_self.subtreeHeight),
      Utils.encodeVarintSigned(_self.subtreeSize),
      Utils.encodeVarintSigned(_self.subtreeVersion),
      uint8(32),  // Size of left subtree hash
      leftSubtree,
      uint8(32),  // Size of right subtree hash
      rightSubtree
    ));
  }
}


/// @title OralceBridge <3 BandChain D3N
/// @author Band Protocol Team
contract OracleBridge {
  using OracleStateMerklePath for OracleStateMerklePath.Data;
  using BlockHeaderMerkleParts for BlockHeaderMerkleParts.Data;
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

/// Mock OracleBridge that allows setting oracle iAVL state at a given height directly.
contract OracleBridgeMock is OracleBridge {
  constructor(address[] memory _validators) public OracleBridge(_validators) {}
  function setOracleState(uint256 _blockHeight, bytes32 _oracleIAVLStateHash) public {
    oracleStates[_blockHeight] = _oracleIAVLStateHash;
  }
}
