pragma solidity 0.6.11;

import {Ownable} from "openzeppelin-solidity/contracts/access/Ownable.sol";
import {IPriceReference} from "./IPriceReference.sol";

contract PriceReferenceProxy is IPriceReference, Ownable {
    IPriceReference public impl;

    constructor(IPriceReference _impl) public {
        patchImpl(_impl);
    }

    function patchImpl(IPriceReference _impl) public onlyOwner {
        impl = _impl;
    }

    function latestRound() public view override returns (uint256) {
        return impl.latestRound();
    }

    function latestAnswer() public view override returns (uint256) {
        return impl.latestAnswer();
    }

    function latestTimestamp() public view override returns (uint256) {
        return impl.latestTimestamp();
    }

    function getAnswer(uint256 _round) public view override returns (uint256) {
        return impl.getAnswer(_round);
    }

    function getTimestamp(uint256 _round) public view override returns (uint256) {
        return impl.getTimestamp(_round);
    }
}
