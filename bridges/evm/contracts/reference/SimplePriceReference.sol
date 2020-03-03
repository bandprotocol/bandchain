pragma solidity 0.5.14;

import {Ownable} from "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import {IPriceReference} from "./IPriceReference.sol";

contract SimplePriceReference is IPriceReference, Ownable {
    uint256[] public prices;
    uint256[] public timestamps;

    constructor(uint256 _initialPrice) public {
        pushPrice(_initialPrice);
    }

    function pushPrice(uint256 _price) public onlyOwner {
        prices.push(_price);
        timestamps.push(block.timestamp);
    }

    function latestRound() public view returns (uint256) {
        return prices.length;
    }

    function latestAnswer() public view returns (uint256) {
        return getAnswer(latestRound());
    }

    function latestTimestamp() public view returns (uint256) {
        return getTimestamp(latestRound());
    }

    function getAnswer(uint256 _round) public view returns (uint256) {
        return prices[_round - 1];
    }

    function getTimestamp(uint256 _round) public view returns (uint256) {
        return timestamps[_round - 1];
    }
}
