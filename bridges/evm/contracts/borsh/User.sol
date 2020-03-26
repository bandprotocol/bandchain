pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import {Borsh} from "./Borsh.sol";
import {ResultDecoder} from "./Result.sol";

contract User {
    using ResultDecoder for bytes;

    function decode(bytes memory _data)
        public
        view
        returns (ResultDecoder.Result memory)
    {
        return _data.decodeResult();
    }
}
