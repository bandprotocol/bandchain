// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {StdReferenceBase} from "./IStdReference.sol";
import {IBridge} from "../interfaces/IBridge.sol";
import {Obi} from "../obi/Obi.sol";

contract StdReference is Ownable, StdReferenceBase {
    using Obi for Obi.Data;

    struct RefData {
        uint64 rate; // USD-rate, multiplied by 1e9.
        uint64 lastUpdate; // UNIX epoch when data is last updated.
    }

    IBridge public bridge; // The underlying bridge smart contract.
    mapping(string => RefData) public refs; // Mapping from symbol to ref data.

    constructor(IBridge _bridge) public {
        bridge = _bridge;
    }

    /// Updates the unlyding bridge implementation. Only callable by the owner.
    function setBridge(IBridge _bridge) public onlyOwner {
        bridge = _bridge;
    }

    /// Calls into the underlying bridge smart contract and saves the result.
    function relayAndVerify(bytes calldata _data) external {
        IBridge.RequestPacket memory req;
        IBridge.ResponsePacket memory res;
        (req, res) = bridge.relayAndVerify(_data);
        _process(req, res);
    }

    /// Calls into the underlying bridge smart contract and saves the result.
    function relayAndMultiVerify(bytes calldata _data) external {
        IBridge.RequestPacket[] memory req;
        IBridge.ResponsePacket[] memory res;
        (req, res) = bridge.relayAndMultiVerify(_data);
        require(req.length == res.length, "INCONSISTENT_PACKET_LENGTH");
        for (uint256 idx = 0; idx < req.length; idx++) {
            _process(req[idx], res[idx]);
        }
    }

    /// Returns the price data for the given base/quote pair. Revert if not available.
    function getReferenceData(string memory _base, string memory _quote)
        public
        override
        view
        returns (ReferenceData memory)
    {
        (uint256 baseRate, uint256 baseLastUpdate) = _getRefData(_base);
        (uint256 quoteRate, uint256 quoteLastUpdate) = _getRefData(_quote);
        return
            ReferenceData({
                rate: (baseRate * 1e18) / quoteRate,
                lastUpdatedBase: baseLastUpdate,
                lastUpdatedQuote: quoteLastUpdate
            });
    }

    function _getRefData(string memory _symbol)
        internal
        view
        returns (uint256 rate, uint256 lastUpdate)
    {
        if (keccak256(bytes(_symbol)) == keccak256(bytes("USD"))) {
            return (1e9, now);
        }
        RefData storage refData = refs[_symbol];
        require(refData.lastUpdate > 0, "REF_DATA_NOT_AVAILABLE");
        return (uint256(refData.rate), uint256(refData.lastUpdate));
    }

    function _process(
        IBridge.RequestPacket memory _req,
        IBridge.ResponsePacket memory _res
    ) internal {
        // TODO: Check request packet
        // TODO: Check response status
        Obi.Data memory reqData = Obi.from(_req.params);
        Obi.Data memory resData = Obi.from(_res.result);
        uint256 symbolLen = uint256(reqData.decodeU32());
        uint256 rateLen = uint256(resData.decodeU32());
        require(symbolLen == rateLen, "INCONSISTENT_REQ_RES");
        for (uint256 idx = 0; idx < symbolLen; idx++) {
            string memory symbol = reqData.decodeString();
            uint64 rate = resData.decodeU64();
            refs[symbol] = RefData(rate, _res.resolveTime);
        }
        uint64 multiplier = reqData.decodeU64();
        require(multiplier == 1e9, "BAD_MULTIPLIER");
    }
}
