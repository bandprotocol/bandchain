// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {IStdReference, StdReferenceBase} from "./IStdReference.sol";

contract StdReferenceProxy is Ownable, StdReferenceBase {
    IStdReference public ref;

    constructor(IStdReference _ref) public {
        ref = _ref;
    }

    /// Updates standard reference implementation. Only callable by the owner.
    function setRef(IStdReference _ref) public onlyOwner {
        ref = _ref;
    }

    /// Returns the price data for the given base/quote pair. Revert if not available.
    function getReferenceData(string memory _base, string memory _quote)
        public
        override
        view
        returns (ReferenceData memory)
    {
        return ref.getReferenceData(_base, _quote);
    }
}
