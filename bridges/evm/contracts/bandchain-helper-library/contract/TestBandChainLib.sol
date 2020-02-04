pragma solidity 0.5.14;

import { BandChainLib } from "./BandChainLib.sol";

contract TestBandChainLib {
  using BandChainLib for bytes;

  function testToUint64List(bytes memory _data) public pure returns(uint64[] memory) {
    return _data.toUint64List();
  }

}
