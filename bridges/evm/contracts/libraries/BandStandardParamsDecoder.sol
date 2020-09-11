pragma solidity ^0.6.0;

import "../obi/Obi.sol";

library ParamsDecoder {
    using Obi for Obi.Data;

    struct Params {
        string symbols;
        uint64 multiplier;
    }

    function decodeParams(bytes memory _data)
        internal
        pure
        returns (Params memory result)
    {
        Obi.Data memory data = Obi.from(_data);
        result.symbols = string(data.decodeBytes());
        result.multiplier = data.decodeU64();
    }
}
