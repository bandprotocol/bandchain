// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;


interface IBridge {
    /// Request packet struct is similar packet on Bandchain using to re-calculate result hash.
    struct RequestPacket {
        string clientId;
        uint64 oracleScriptId;
        bytes params;
        uint64 askCount;
        uint64 minCount;
    }

    /// Response packet struct is similar packet on Bandchain using to re-calculate result hash.
    struct ResponsePacket {
        string clientId;
        uint64 requestId;
        uint64 ansCount;
        uint64 requestTime;
        uint64 resolveTime;
        uint8 resolveStatus;
        bytes result;
    }

    struct ValidatorWithPower {
        address addr;
        uint256 power;
    }

    /// Update validator powers by owner.
    /// @param _validators The changed set of BandChain validators.
    function updateValidatorPowers(ValidatorWithPower[] calldata _validators) external;

    /// Performs oracle state relay and oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relayAndVerify(bytes calldata _data)
        external
        returns (RequestPacket memory, ResponsePacket memory);
}
