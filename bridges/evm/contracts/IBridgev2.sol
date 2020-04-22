pragma solidity 0.5.14;

import {IBridge} from "./IBridge.sol";


interface IBridegeV2 {
    /// Event to broadcast oracle request event.
    event OracleRequest(
        string clientId,
        uint64 oracleScriptId,
        bytes params,
        uint64 askCount,
        uint64 minCount
    );

    /// Requests a new oracle script to bandchain by emit request event.
    function requestOracle(
        uint64 oracleScriptId,
        bytes calldata params,
        uint64 askCount,
        uint64 minCount
    ) external returns (bytes32);

    /// Performs oracle state relay, oracle data verification and contract function calling
    /// The caller submits the encoded proof and target contract will receive the decoded data,
    /// ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relayAndCall(bytes calldata _data) external;
}
