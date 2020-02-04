pragma solidity 0.5.14;

import { BandChainLib } from "../BandChainLib.sol";

contract BandChainLibMock {
  using BandChainLib for bytes;

  function toUint64List(bytes memory _data) public pure returns(uint64[] memory) {
    return _data.toUint64List();
  }

}
