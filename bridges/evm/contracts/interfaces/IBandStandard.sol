// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

interface IBandStandard {
    struct SymbolData {
        uint64 oracleScriptID;
        uint8 calldataID;
        uint8 rateID;
    }
    
    struct DataUpdate {
        uint256 base;
        uint256 quote;
    }
    
    struct ReferenceData {
        uint256 rate;
        DataUpdate lastUpdated;
    }
    
    function getReferenceData(string[] memory pairs)
    external
    view
    returns (ReferenceData[] memory);
}
