pragma solidity ^0.5.0;

import "./Borsh.sol";

library ResultDecoder {
    using Borsh for Borsh.Data;

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
        Borsh.Data memory data = Borsh.from(_data);
        result.symbol = string(data.decodeBytes());
        result.multiplier = data.decodeU64();
        result.what = data.decodeU8();
    }
}
