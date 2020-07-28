// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;

import {BandChainLib} from "../BandChainLib.sol";

contract BandChainLibMock {
    using BandChainLib for bytes;

    function toUint64List(bytes memory _data)
        public
        pure
        returns (uint64[] memory)
    {
        return _data.toUint64List();
    }

}
