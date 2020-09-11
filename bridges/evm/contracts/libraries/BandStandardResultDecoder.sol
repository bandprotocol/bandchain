pragma solidity ^0.6.0;

import "../obi/Obi.sol";

library ResultDecoder {
    using Obi for Obi.Data;

    struct Result {
        uint64[] rates;
    }

    function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
    {
        Obi.Data memory data = Obi.from(_data);
        uint32 length = data.decodeU32();
        uint64[] memory rates = new uint64[](length);
        for (uint256 i = 0; i < length; i++) {
            rates[i] = data.decodeU64();
        }
        return Result(rates);
    }
}
