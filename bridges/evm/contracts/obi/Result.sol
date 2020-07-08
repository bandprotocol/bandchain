// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;

import "./Obi.sol";

library ResultDecoder {
    using Obi for Obi.Data;

    struct Result {
        string symbol;
        uint64 multiplier;
        uint8 what;
    }

    function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
    {
        Obi.Data memory data = Obi.from(_data);
        result.symbol = string(data.decodeBytes());
        result.multiplier = data.decodeU64();
        result.what = data.decodeU8();
    }
}
