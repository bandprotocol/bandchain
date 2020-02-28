pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;
import {BlockHeaderMerkleParts} from "./BlockHeaderMerkleParts.sol";
import {IAVLMerklePath} from "./IAVLMerklePath.sol";
import {TMSignature} from "./TMSignature.sol";
import {Utils} from "./Utils.sol";
import {IBridge} from "./IBridge.sol";

/// @title Bridge <3 BandChain D3N
/// @author Band Protocol Team
contract Bridge is IBridge {
    using BlockHeaderMerkleParts for BlockHeaderMerkleParts.Data;
    using IAVLMerklePath for IAVLMerklePath.Data;
    using TMSignature for TMSignature.Data;

    /// Mapping from block height to the hash of "zoracle" iAVL Merkle tree.
    mapping(uint256 => bytes32) public oracleStates;
    /// Mapping from an address to whether it's a validator.
    mapping(address => bool) public validators;
    /// The total number of active validators currently on duty.
    uint256 public validatorCount;

    /// Initializes an oracle bridge to BandChain.
    /// @param _validators The initial set of BandChain active validators.
    constructor(address[] memory _validators) public {
        for (uint256 idx = 0; idx < _validators.length; ++idx) {
            address validator = _validators[idx];
            require(
                !validators[validator],
                "DUPLICATION_IN_INITIAL_VALIDATOR_SET"
            );
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
        bytes32 _supplyStoresMerkleHash,
        BlockHeaderMerkleParts.Data memory _merkleParts,
        bytes memory _signedDataPrefix,
        TMSignature.Data[] memory _signatures
    ) public {
        // Computes Tendermint's application state hash at this given block. AppHash is actually a
        // Merkle hash on muliple stores. Luckily, we only care about "zoracle" tree and all other
        // stores can just be combined into one bytes32 hash off-chain.
        //
        //                                            ____________appHash_________
        //                                          /                              \
        //                   ____otherStoresMerkleHash ____                         ___innerHash___
        //                 /                                \                     /                  \
        //         _____ h5 ______                    ______ h6 _______        supply              zoracle
        //       /                \                 /                  \
        //     h1                  h2             h3                    h4
        //     /\                  /\             /\                    /\
        //  acc  distribution   gov  main     mint  params     slashing   staking
        bytes32 appHash = Utils.merkleInnerHash(
            _otherStoresMerkleHash,
            Utils.merkleInnerHash(
                _supplyStoresMerkleHash,
                Utils.merkleLeafHash(
                    abi.encodePacked(
                        hex"077a6f7261636c6520", // uint8(7) + "zoracle" + uint8(32)
                        sha256(
                            abi.encodePacked(
                                sha256(abi.encodePacked(_oracleIAVLStateHash))
                            )
                        )
                    )
                )
            )
        );
        // Computes Tendermint's block header hash at this given block.
        bytes32 blockHeader = _merkleParts.getBlockHeader(
            appHash,
            _blockHeight
        );
        // Counts the total number of valid signatures signed by active validators.
        address lastSigner = address(0);
        uint256 validSignatureCount = 0;
        for (uint256 idx = 0; idx < _signatures.length; ++idx) {
            address signer = _signatures[idx].recoverSigner(
                blockHeader,
                _signedDataPrefix
            );
            require(signer > lastSigner, "INVALID_SIGNATURE_SIGNER_ORDER");
            if (validators[signer]) {
                // Increases valid signature count if the signer is one of the validators.
                validSignatureCount += 1;
            }
            lastSigner = signer;
        }
        // Verifies that sufficient validators signed the block and saves the oracle state.
        require(
            validSignatureCount * 3 > validatorCount * 2,
            "INSUFFICIENT_VALIDATOR_SIGNATURES"
        );
        oracleStates[_blockHeight] = _oracleIAVLStateHash;
    }

    /// Helper struct to workaround Solidity's "stack too deep" problem.
    struct VerifyOracleDataLocalVariables {
        bytes encodedVarint;
        bytes32 dataHash;
    }

    /// Verifies that the given data is a valid data on BandChain as of the given block height.
    /// @param _blockHeight The block height. Someone must already relay this block.
    /// @param _data The data to verify, with the format similar to what on the blockchain store.
    /// @param _requestId The ID of request for this data piece.
    /// @param _version Lastest block height that the data node was updated.
    /// @param _merklePaths Merkle proof that shows how the data leave is part of the oracle iAVL.
    function verifyOracleData(
        uint256 _blockHeight,
        bytes memory _data,
        uint64 _requestId,
        uint64 _oracleScriptId,
        bytes memory _params,
        uint256 _version,
        IAVLMerklePath.Data[] memory _merklePaths
    ) public view returns (VerifyOracleDataResult memory) {
        bytes32 oracleStateRoot = oracleStates[_blockHeight];
        require(
            oracleStateRoot != bytes32(uint256(0)),
            "NO_ORACLE_ROOT_STATE_DATA"
        );
        // Computes the hash of leaf node for iAVL oracle tree.
        VerifyOracleDataLocalVariables memory vars;
        vars.encodedVarint = Utils.encodeVarintSigned(_version);
        vars.dataHash = sha256(_data);
        bytes32 currentMerkleHash = sha256(
            abi.encodePacked(
                uint8(0), // Height of tree (only leaf node) is 0 (signed-varint encode)
                uint8(2), // Size of subtree is 1 (signed-varint encode)
                vars.encodedVarint,
                uint8(17 + _params.length), // Size of data key (1-byte constant 0x01 + 8-byte request ID + 8-byte oracleScriptId + length of params)
                uint8(255), // Constant 0xff prefix data request info storage key
                _requestId,
                _oracleScriptId,
                _params,
                uint8(32), // Size of data hash
                vars.dataHash
            )
        );
        // Goes step-by-step computing hash of parent nodes until reaching root node.
        for (uint256 idx = 0; idx < _merklePaths.length; ++idx) {
            currentMerkleHash = _merklePaths[idx].getParentHash(
                currentMerkleHash
            );
        }
        // Verifies that the computed Merkle root matches what currently exists.
        require(
            currentMerkleHash == oracleStateRoot,
            "INVALID_ORACLE_DATA_PROOF"
        );
        return VerifyOracleDataResult(_data, _oracleScriptId, _params);
    }
    /// Decode the encoded result and receives back the decoded data which is the data and its context.
    /// @param _data The encoded of result and its context.
    function decodeWrappedResult(bytes memory _data)
        public
        pure
        returns (WrappedResult memory)
    {
        require(_data.length > 40, "INPUT_MUST_BE_LONGER_THAN_40_BYTES");

        WrappedResult memory wrapped;
        assembly {
            mstore(
                add(wrapped, 0x00),
                and(mload(add(_data, add(0x08, 0x00))), 0xffffffffffffffff)
            )
            mstore(
                add(wrapped, 0x20),
                and(mload(add(_data, add(0x08, 0x08))), 0xffffffffffffffff)
            )
            mstore(
                add(wrapped, 0x40),
                and(mload(add(_data, add(0x08, 0x10))), 0xffffffffffffffff)
            )
            mstore(
                add(wrapped, 0x60),
                and(mload(add(_data, add(0x08, 0x18))), 0xffffffffffffffff)
            )
            mstore(
                add(wrapped, 0x80),
                and(mload(add(_data, add(0x08, 0x20))), 0xffffffffffffffff)
            )
        }

        bytes memory data = new bytes(_data.length - 40);
        uint256 l = ((data.length - 1) / 32) + 1;
        for (uint256 i = 0; i < l; i++) {
            assembly {
                mstore(
                    add(data, add(0x20, mul(i, 0x20))),
                    mload(add(_data, add(0x48, mul(i, 0x20))))
                )
            }
        }
        wrapped.data = data;

        return wrapped;
    }

    /// Verifies that the given data is a valid data on BandChain as of the given block height.
    /// @param _blockHeight The block height. Someone must already relay this block.
    /// @param _data The data to verify, with the format similar to what on the blockchain store.
    /// @param _requestId The ID of request for this data piece.
    /// @param _version Lastest block height that the data node was updated.
    /// @param _merklePaths Merkle proof that shows how the data leave is part of the oracle iAVL.
    function verifyOracleData(
        uint256 _blockHeight,
        bytes memory _data,
        uint64 _requestId,
        uint64 _oracleScriptId,
        bytes memory _params,
        uint256 _version,
        IAVLMerklePath.Data[] memory _merklePaths
    ) public view returns (VerifyOracleDataResult memory) {
        bytes32 oracleStateRoot = oracleStates[_blockHeight];
        require(
            oracleStateRoot != bytes32(uint256(0)),
            "NO_ORACLE_ROOT_STATE_DATA"
        );
        // Computes the hash of leaf node for iAVL oracle tree.
        VerifyOracleDataLocalVariables memory vars;
        vars.encodedVarint = Utils.encodeVarintSigned(_version);
        vars.dataHash = sha256(_data);
        bytes32 currentMerkleHash = sha256(
            abi.encodePacked(
                uint8(0), // Height of tree (only leaf node) is 0 (signed-varint encode)
                uint8(2), // Size of subtree is 1 (signed-varint encode)
                vars.encodedVarint,
                uint8(17 + _params.length), // Size of data key (1-byte constant 0x01 + 8-byte request ID + 8-byte oracleScriptId + length of params)
                uint8(255), // Constant 0xff prefix data request info storage key
                _requestId,
                _oracleScriptId,
                _params,
                uint8(32), // Size of data hash
                vars.dataHash
            )
        );
        // Goes step-by-step computing hash of parent nodes until reaching root node.
        for (uint256 idx = 0; idx < _merklePaths.length; ++idx) {
            currentMerkleHash = _merklePaths[idx].getParentHash(
                currentMerkleHash
            );
        }
        // Verifies that the computed Merkle root matches what currently exists.
        require(
            currentMerkleHash == oracleStateRoot,
            "INVALID_ORACLE_DATA_PROOF"
        );
        return
            VerifyOracleDataResult(
                decodeWrappedResult(_data),
                _oracleScriptId,
                _params
            );
    }

    /// Performs oracle state relay and oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relayAndVerify(bytes calldata _data)
        external
        returns (VerifyOracleDataResult memory result)
    {
        (bytes memory relayData, bytes memory verifyData) = abi.decode(
            _data,
            (bytes, bytes)
        );
        (bool relayOk, ) = address(this).call(
            abi.encodePacked(this.relayOracleState.selector, relayData)
        );
        require(relayOk, "RELAY_ORACLE_STATE_FAILED");
        (bool verifyOk, bytes memory verifyResult) = address(this).staticcall(
            abi.encodePacked(this.verifyOracleData.selector, verifyData)
        );
        require(verifyOk, "VERIFY_ORACLE_DATA_FAILED");
        return abi.decode(verifyResult, (VerifyOracleDataResult));
    }
}
