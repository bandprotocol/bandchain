// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {StdReferenceBase} from "./IStdReference.sol";
import {IBridge} from "../interfaces/IBridge.sol";
import {Obi} from "../obi/Obi.sol";

contract OptimisticStdReference {

}

// contract OptimisticStdReference is Ownable, StdReferenceBase {
//     struct RefData {
//         uint64 rate; // USD-rate, multiplied by 1e9.
//         uint64 resolveTime; // UNIX epoch when data is resolved.
//         uint64 requestId; // Request ID that contains this ref data.
//         uint64 relayTime; // UNIX epoch when data is relayed to this contract.
//     }

//     mapping(string => RefData) public pending;
//     mapping(string => RefData) public confirmed;

//     function relay(
//         string[] calldata _symbols,
//         uint64[] calldata _rates,
//         uint64[] calldata _resolveTimes,
//         uint64[] calldata _requestIds
//     ) external onlyOwner {
//         uint256 len = _symbols.length;
//         require(_rates.length == len, "BAD_RATES_LENGTH");
//         require(_resolveTimes.length == len, "BAD_RESOLVE_TIMES_LENGTH");
//         require(_requestIds.length == len, "BAD_REQUEST_IDS_LENGTH");
//         for (uint256 idx = 0; idx < len; idx++) {
//             _relay(
//                 _symbols[idx],
//                 _rates[idx],
//                 _resolveTimes[idx],
//                 _requestIds[idx]
//             );
//         }
//     }

//     function _relay(
//         string memory _symbol,
//         uint64 _rate,
//         uint64 _resolveTime,
//         uint64 _requestId
//     ) internal {
//         // RefData storage pendingRef = pending[_symbol];
//         // confirmed[_symbol] = pending[_symbol];
//     }

//     function getReferenceData(string memory _base, string memory _quote)
//         public
//         override
//         view
//         returns (ReferenceData memory)
//     {
//         (uint256 baseRate, uint256 baseLastUpdate) = _getRefData(_base);
//         (uint256 quoteRate, uint256 quoteLastUpdate) = _getRefData(_quote);
//         return
//             ReferenceData({
//                 rate: (baseRate * 1e18) / quoteRate,
//                 lastUpdatedBase: baseLastUpdate,
//                 lastUpdatedQuote: quoteLastUpdate
//             });
//     }

//     function _getRefData(string memory _symbol)
//         internal
//         view
//         returns (uint256 rate, uint256 lastUpdate)
//     {
//         if (keccak256(bytes(_symbol)) == keccak256(bytes("USD"))) {
//             return (1e9, now);
//         }
//         RefData storage refData = confirmed[_symbol];
//         require(refData.resolveTime > 0, "REF_DATA_NOT_AVAILABLE");
//         return (uint256(refData.rate), uint256(refData.resolveTime));
//     }
// }
