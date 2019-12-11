pragma solidity 0.5.14;

library UintLib {

  function toU8Arr(uint256 prefix) internal pure returns(bytes memory) {
    uint256 n = (prefix >> 248) & 127;
    bytes memory arr = new bytes(n);
    while (n > 0) {
      arr[n-1] = byte(uint8(prefix & 255));
      prefix >>= 8;
      n--;
    }
    return arr;
  }

}
