// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {IBandStandard} from "./interfaces/IBandStandard.sol";
import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";

contract BandOracleAggregatorProxy is IBandStandard, Ownable {
    
    IBandStandard public aggregator;
    
    constructor(IBandStandard _aggregator) public {
        aggregator = _aggregator;
    }
    
    function setAggregator(IBandStandard _aggregator) public onlyOwner {
        aggregator = _aggregator;
    }
    
    function getReferenceData(string[] memory pairs)
    external
    override
    view
    returns (ReferenceData[] memory) {
        return aggregator.getReferenceData(pairs);
    }
}
