pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {Packets} from "./Packets.sol";
import {Bridge} from "./Bridge.sol";
import {IBridgeCache} from "./IBridgeCache.sol";

contract BridgeProxy is IBridge, IBridgeCache, Ownable {
    IBridgeCache public bridge;

    /**
     * @notice Contract constructor
     * @dev Initializes a new BandPriceOracle instance.
     * @param _bridge Band's Bridge contract address
     **/
    constructor(IBridgeCache _bridge) public {
        bridge = _bridge;
    }

    /// Set the address of the bridge contract to use
    /// @param _bridge The address of the bridge to use
    function setBridge(IBridgeCache _bridge) external onlyOwner {
        bridge = _bridge;
    }
    
    /// Returns the hash of a RequestPacket.
    /// @param _request A tuple that represents RequestPacket struct.
    function getRequestKey(RequestPacket memory _request)
        external
        view
        override
        returns (bytes32)
    {
        return bridge.getRequestKey(_request);
    }
    
    /// Returns the ResponsePacket for a given RequestPacket.
    /// Reverts if can't find the related response in the mapping.
    /// @param _request A tuple that represents RequestPacket struct.
    function getLatestResponse(RequestPacket memory _request)
        external
        view
        override
        returns (ResponsePacket memory)
    {
        return bridge.getLatestResponse(_request);
    }

    /// Update validator powers by owner.
    /// @param _validators The changed set of BandChain validators.
    function updateValidatorPowers(ValidatorWithPower[] memory _validators)
        external
        override
        onlyOwner
    {
      bridge.updateValidatorPowers(_validators);
    }

    function relayAndVerify(bytes calldata _data)
        external
        override
        returns (RequestPacket memory, ResponsePacket memory)
    {
        return bridge.relayAndVerify(_data);
    }

    function relay(bytes calldata _data)
        external
        override
    {
        return bridge.relay(_data);
    }
}

