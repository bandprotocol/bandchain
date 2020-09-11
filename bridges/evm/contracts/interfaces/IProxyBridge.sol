// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {IBridge} from "./IBridge.sol";
import {ICacheBridge} from "./ICacheBridge.sol";

interface IProxyBridge is IBridge, ICacheBridge {

    /// Set the address of the bridge contract to use
    /// @param _bridge The address of the bridge to use
    function setBridge(ICacheBridge _bridge) external;
    
    /// Returns the ResponsePacket for a given RequestPacket.
    /// Reverts if can't find the related response in the mapping.
    /// @param _request A tuple that represents RequestPacket struct.
    function getLatestResponse(IBridge.RequestPacket memory _request) external override view returns (IBridge.ResponsePacket memory);


    /// Performs oracle state relay and oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relayAndVerify(bytes calldata _data) external override returns (IBridge.RequestPacket memory, IBridge.ResponsePacket memory);

    /// Performs oracle state relay and many times of oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and an array of data verification.
    function relayAndMultiVerify(bytes calldata _data)
        external
        override
        returns (RequestPacket[] memory, ResponsePacket[] memory);

    /// Performs oracle state relay and oracle data verification in one go.
    /// After that, the results will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relay(bytes calldata _data) external override;

    /// Performs oracle state relay and many times of oracle data verification in one go.
    /// After that, the results which is an array of Packet will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and an array of data verification.
    function relayMulti(bytes calldata _data) external override;
}

