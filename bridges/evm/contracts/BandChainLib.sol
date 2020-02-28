pragma solidity 0.5.14;

library BandChainLib {
    function toUint64List(bytes memory _data)
        internal
        pure
        returns (uint64[] memory)
    {
        uint64[] memory result = new uint64[](_data.length / 8);
        require(result.length * 8 == _data.length, "DATA_LENGTH_IS_INVALID");

    struct Result {
      uint64 requestTime;
			uint64 aggregationTime;
			uint64 requestedValidatorsCount;
			uint64 sufficientValidatorCount;
			uint64 reportedValidatorsCount;
			bytes data;
    }

  function decodeToResult(bytes memory _data) public pure returns (Result memory) {
    require(_data.length > 40, "INPUT_MUST_BE_LONGER_THAN_40_BYTES");

    Result memory result;
    assembly {
      mstore(add(result, 0x00), and(mload(add(_data, add(0x08, 0x00))), 0xffffffffffffffff))
      mstore(add(result, 0x20), and(mload(add(_data, add(0x08, 0x08))), 0xffffffffffffffff))
      mstore(add(result, 0x40), and(mload(add(_data, add(0x08, 0x10))), 0xffffffffffffffff))
      mstore(add(result, 0x60), and(mload(add(_data, add(0x08, 0x18))), 0xffffffffffffffff))
      mstore(add(result, 0x80), and(mload(add(_data, add(0x08, 0x20))), 0xffffffffffffffff))
    }

    bytes memory data = new bytes(_data.length - 40);
    uint256 l = ((data.length - 1) / 32) + 1;
    for (uint256 i = 0; i < l; i++) {
      assembly {
        mstore(add(data,add(0x20, mul(i,0x20))), mload(add(_data, add(0x48, mul(i,0x20)))))
      }
    }
    result.data = data;

    return  result;
  }

  function toUint64List(bytes memory _data) internal pure returns(uint64[] memory) {
    uint64[] memory result = new uint64[](_data.length / 8);
    require(result.length * 8 == _data.length, "DATA_LENGTH_IS_INVALID");

        for (uint256 i = 0; i < result.length; i++) {
            uint64 tmp;
            assembly {
                tmp := mload(add(_data, add(0x08, mul(i, 0x08))))
            }
            result[i] = tmp;
        }

        return result;
    }
}
