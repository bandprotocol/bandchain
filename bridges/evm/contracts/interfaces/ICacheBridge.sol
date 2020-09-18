// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {IBridge} from "./IBridge.sol";

interface ICacheBridge is IBridge {
    /// Returns the ResponsePacket for a given RequestPacket.
    /// @param _request The tuple that represents RequestPacket struct.
    function getLatestResponse(RequestPacket calldata _request)
        external
        view
        returns (ResponsePacket memory);

    /// Performs oracle state relay and oracle data verification in one go.
    /// After that, the results will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relay(bytes calldata _data) external;

    /// Performs oracle state relay and many times of oracle data verification in one go.
    /// After that, the results which is an array of Packet will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and an array of data verification.
    function relayMulti(bytes calldata _data) external;
}
