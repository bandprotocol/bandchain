pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;


interface IBridge {
    /// Request packet struct is similar packet on Bandchain using to re-calculate result hash.
    struct RequestPacket {
        string clientId;
        uint64 oracleScriptId;
        string params;
        uint64 askCount;
        uint64 minCount;
    }

    /// Response packet struct is similar packet on Bandchain using to re-calculate result hash.
    struct ResponsePacket {
        string clientId;
        uint64 requestId;
        uint64 ansCount;
        uint64 prepareTime;
        uint64 resolveTime;
        uint8 resolveStatus;
        string data;
    }

    /// Helper struct to help the function caller to decode oracle data.
    struct VerifyOracleDataResult {
        uint64 oracleScriptId;
        uint64 prepareTime;
        uint64 resolveTime;
        uint64 askCount;
        uint64 minCount;
        uint64 ansCount;
        bytes params;
        bytes data;
    }

    /// Performs oracle state relay and oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relayAndVerify(bytes calldata _data)
        external
        returns (VerifyOracleDataResult memory result);
}
