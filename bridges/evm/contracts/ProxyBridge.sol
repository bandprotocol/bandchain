// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {IProxyBridge} from "./interfaces/IProxyBridge.sol";
import {ICacheBridge} from "./interfaces/ICacheBridge.sol";

contract ProxyBridge is IProxyBridge, Ownable {
    ICacheBridge public bridge;

    /// @notice Contract constructor
    /// @dev Initializes a new Bridge instance.
    /// @param _bridge Band's Bridge contract address
    constructor(ICacheBridge _bridge) public {
        bridge = _bridge;
    }

    /// Set the address of the bridge contract to use
    /// @param _bridge The address of the bridge to use
    function setBridge(ICacheBridge _bridge) external override onlyOwner {
        bridge = _bridge;
    }

    /// Returns the ResponsePacket for a given RequestPacket.
    /// Reverts if can't find the related response in the mapping.
    /// @param _request A tuple that represents RequestPacket struct.
    function getLatestResponse(ICacheBridge.RequestPacket memory _request)
        external
        override
        view
        returns (ICacheBridge.ResponsePacket memory)
    {
        return bridge.getLatestResponse(_request);
    }

    /// Performs oracle state relay and oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relayAndVerify(bytes calldata _data)
        external
        override
        returns (
            ICacheBridge.RequestPacket memory,
            ICacheBridge.ResponsePacket memory
        )
    {
        return bridge.relayAndVerify(_data);
    }

    /// Performs oracle state relay and many times of oracle data verification in one go. The caller submits
    /// the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param _data The encoded data for oracle state relay and an array of data verification.
    function relayAndMultiVerify(bytes calldata _data)
        external
        override
        returns (
            ICacheBridge.RequestPacket[] memory,
            ICacheBridge.ResponsePacket[] memory
        )
    {
        return bridge.relayAndMultiVerify(_data);
    }

    /// Performs oracle state relay and oracle data verification in one go.
    /// After that, the results will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and data verification.
    function relay(bytes calldata _data) external override {
        return bridge.relay(_data);
    }

    /// Performs oracle state relay and many times of oracle data verification in one go.
    /// After that, the results which is an array of Packet will be recorded to the state by using the hash of RequestPacket as key.
    /// @param _data The encoded data for oracle state relay and an array of data verification.
    function relayMulti(bytes calldata _data) external override {
        return bridge.relayMulti(_data);
    }
}
