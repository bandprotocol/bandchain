// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {AccessControl} from "openzeppelin-solidity/contracts/access/AccessControl.sol"; // prettier-ignore
import {StdReferenceBase} from "./IStdReference.sol";

contract StdReferenceBasic is AccessControl, StdReferenceBase {
    event RefDataUpdate(
        string symbol,
        uint64 rate,
        uint64 resolveTime,
        uint64 requestId
    );

    struct RefData {
        uint64 rate; // USD-rate, multiplied by 1e9.
        uint64 resolveTime; // UNIX epoch when data is last resolved.
        uint64 requestId; // BandChain request identifier for this data.
    }

    mapping(string => RefData) public refs; // Mapping from symbol to ref data.
    bytes32 public constant RELAYER_ROLE = keccak256("RELAYER_ROLE");

    constructor() public {
        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _setupRole(RELAYER_ROLE, msg.sender);
    }

    function relay(
        string[] memory _symbols,
        uint64[] memory _rates,
        uint64[] memory _resolveTimes,
        uint64[] memory _requestIds
    ) external {
        require(hasRole(RELAYER_ROLE, msg.sender), "NOT_A_RELAYER");
        uint256 len = _symbols.length;
        require(_rates.length == len, "BAD_RATES_LENGTH");
        require(_resolveTimes.length == len, "BAD_RESOLVE_TIMES_LENGTH");
        require(_requestIds.length == len, "BAD_REQUEST_IDS_LENGTH");
        for (uint256 idx = 0; idx < len; idx++) {
            refs[_symbols[idx]] = RefData({
                rate: _rates[idx],
                resolveTime: _resolveTimes[idx],
                requestId: _requestIds[idx]
            });
            emit RefDataUpdate(
                _symbols[idx],
                _rates[idx],
                _resolveTimes[idx],
                _requestIds[idx]
            );
        }
    }

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
        require(refData.resolveTime > 0, "REF_DATA_NOT_AVAILABLE");
        return (uint256(refData.rate), uint256(refData.resolveTime));
    }
}
