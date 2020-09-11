// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {IBandDataset} from "./interfaces/IBandDataset.sol";
import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";

contract BandDatasetProxy is IBandDataset, Ownable {
    IBandDataset public dataset;

    constructor(IBandDataset _dataset) public {
        dataset = _dataset;
    }

    function setAggregator(IBandDataset _dataset) public onlyOwner {
        dataset = _dataset;
    }

    function getReferenceData(string[] memory pairs)
        external
        override
        view
        returns (ReferenceData[] memory)
    {
        return dataset.getReferenceData(pairs);
    }
}
