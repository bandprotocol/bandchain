// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;
import {BlockHeaderMerkleParts} from "./BlockHeaderMerkleParts.sol";
import {MultiStore} from "./MultiStore.sol";
import {SafeMath} from "openzeppelin-solidity/contracts/math/SafeMath.sol";
import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {IAVLMerklePath} from "./IAVLMerklePath.sol";
import {TMSignature} from "./TMSignature.sol";
import {Utils} from "./Utils.sol";
import {Packets} from "./Packets.sol";
import {IBridge} from "./IBridge.sol";

/// @title Bridge <3 BandChain D3N
/// @author Band Protocol Team
contract Bridge is IBridge, Ownable {
    using BlockHeaderMerkleParts for BlockHeaderMerkleParts.Data;
    using MultiStore for MultiStore.Data;
    using IAVLMerklePath for IAVLMerklePath.Data;
    using TMSignature for TMSignature.Data;
    using SafeMath for uint256;

    struct ValidatorWithPower {
        address addr;
        uint256 power;
    }

    /// Mapping from block height to the hash of "zoracle" iAVL Merkle tree.
    mapping(uint256 => bytes32) public oracleStates;
    /// Mapping from an address to its voting power.
    mapping(address => uint256) public validatorPowers;
    /// The total voting power of active validators currently on duty.
    uint256 public totalValidatorPower;

    /// Initializes an oracle bridge to BandChain.
    /// @param _validators The initial set of BandChain active validators.
    constructor(ValidatorWithPower[] memory _validators) public {
        for (uint256 idx = 0; idx < _validators.length; ++idx) {
            ValidatorWithPower memory validator = _validators[idx];
            require(
                validatorPowers[validator.addr] == 0,
                "DUPLICATION_IN_INITIAL_VALIDATOR_SET"
            );
            validatorPowers[validator.addr] = validator.power;
            totalValidatorPower = totalValidatorPower.add(validator.power);
        }
    }

    /// Update validator powers by owner.
    /// @param _validators The changed set of BandChain validators.
    function updateValidatorPowers(ValidatorWithPower[] memory _validators)
        external
        onlyOwner
    {
        for (uint256 idx = 0; idx < _validators.length; ++idx) {
            ValidatorWithPower memory validator = _validators[idx];
            totalValidatorPower = totalValidatorPower.sub(
                validatorPowers[validator.addr]
            );
            validatorPowers[validator.addr] = validator.power;
            totalValidatorPower = totalValidatorPower.add(validator.power);
        }
    }

    /// Relays a new oracle state to the bridge contract.
    /// @param _blockHeight The height of block to relay to this bridge contract.
    /// @param _multiStore Extra multi store to compute app hash. See MultiStore lib.
    /// @param _merkleParts Extra merkle parts to compute block hash. See BlockHeaderMerkleParts lib.
    /// @param _signatures The signatures signed on this block, sorted alphabetically by address.
    function relayOracleState(
        uint256 _blockHeight,
        MultiStore.Data memory _multiStore,
        BlockHeaderMerkleParts.Data memory _merkleParts,
        TMSignature.Data[] memory _signatures
    ) public {
        bytes32 appHash = _multiStore.getAppHash();
        // Computes Tendermint's block header hash at this given block.
        bytes32 blockHeader = _merkleParts.getBlockHeader(
            appHash,
            _blockHeight
        );
        // Counts the total number of valid signatures signed by active validators.
        address lastSigner = address(0);
        uint256 sumVotingPower = 0;
        for (uint256 idx = 0; idx < _signatures.length; ++idx) {
            address signer = _signatures[idx].recoverSigner(blockHeader);
            require(signer > lastSigner, "INVALID_SIGNATURE_SIGNER_ORDER");
            sumVotingPower = sumVotingPower.add(validatorPowers[signer]);
            lastSigner = signer;
        }
        // Verifies that sufficient validators signed the block and saves the oracle state.
        require(
            sumVotingPower.mul(3) > totalValidatorPower.mul(2),
            "INSUFFICIENT_VALIDATOR_SIGNATURES"
        );
        oracleStates[_blockHeight] = _multiStore.oracleIAVLStateHash;
    }

    /// Helper struct to workaround Solidity's "stack too deep" problem.
    struct VerifyOracleDataLocalVariables {
        bytes encodedVarint;
        bytes32 dataHash;
    }

    /// Verifies that the given data is a valid data on BandChain as of the given block height.
    /// @param _blockHeight The block height. Someone must already relay this block.
    /// @param _requestPacket The request packet is this request.
    /// @param _responsePacket The response packet of this request.
    /// @param _version Lastest block height that the data node was updated.
    /// @param _merklePaths Merkle proof that shows how the data leave is part of the oracle iAVL.
    function verifyOracleData(
        uint256 _blockHeight,
        RequestPacket memory _requestPacket,
        ResponsePacket memory _responsePacket,
        uint256 _version,
        IAVLMerklePath.Data[] memory _merklePaths
    ) public view returns (RequestPacket memory, ResponsePacket memory) {
        bytes32 oracleStateRoot = oracleStates[_blockHeight];
        require(
            oracleStateRoot != bytes32(uint256(0)),
            "NO_ORACLE_ROOT_STATE_DATA"
        );
        // Computes the hash of leaf node for iAVL oracle tree.
        VerifyOracleDataLocalVariables memory vars;
        vars.encodedVarint = Utils.encodeVarintSigned(_version);
        vars.dataHash = sha256(
            Packets.getEncodedResult(_requestPacket, _responsePacket)
        );
        bytes32 currentMerkleHash = sha256(
            abi.encodePacked(
                uint8(0), // Height of tree (only leaf node) is 0 (signed-varint encode)
                uint8(2), // Size of subtree is 1 (signed-varint encode)
                vars.encodedVarint,
                uint8(9), // Size of data key (1-byte constant 0x01 + 8-byte request ID)
                uint8(255), // Constant 0xff prefix data request info storage key
                _responsePacket.requestId,
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

        return (_requestPacket, _responsePacket);
    }

    /// Performs oracle state relay and oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relayAndVerify(bytes calldata _data)
        external
        override
        returns (RequestPacket memory, ResponsePacket memory)
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
        return abi.decode(verifyResult, (RequestPacket, ResponsePacket));
    }

    /// Performs oracle state relay and many times of oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and an array of data verification.
    function relayAndMultiVerify(bytes calldata _data)
        external
        override
        returns (Packet[] memory)
    {
        (bytes memory relayData, bytes[] memory manyVerifyData) = abi.decode(
            _data,
            (bytes, bytes[])
        );
        (bool relayOk, ) = address(this).call(
            abi.encodePacked(this.relayOracleState.selector, relayData)
        );
        require(relayOk, "RELAY_ORACLE_STATE_FAILED");

        Packet[] memory packets = new Packet[](manyVerifyData.length);
        for (uint256 i = 0; i < manyVerifyData.length; i++) {
            (bool verifyOk, bytes memory verifyResult) = address(this)
                .staticcall(
                abi.encodePacked(
                    this.verifyOracleData.selector,
                    manyVerifyData[i]
                )
            );
            require(verifyOk, "VERIFY_ORACLE_DATA_FAILED");
            (packets[i].request, packets[i].response) = abi.decode(
                verifyResult,
                (RequestPacket, ResponsePacket)
            );
        }

        return packets;
    }
}
