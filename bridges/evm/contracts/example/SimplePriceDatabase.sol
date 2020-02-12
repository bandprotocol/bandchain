pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import { BandChainLib } from "../BandChainLib.sol";
import { IBridge } from "../IBridge.sol";

contract SimplePriceDatabase {
  using BandChainLib for bytes;

  bytes32 public codeHash;
  bytes public params;
  uint256 public latestETHPrice;
  uint256 public lastUpdate;

  IBridge public bridge;

  constructor(bytes32 _codeHash , bytes memory _params, IBridge _bridge) public {
    codeHash = _codeHash;
    params = _params;
    bridge = _bridge;
  }

  function update(bytes memory _reportPrice) public {
    IBridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(_reportPrice);

    require(result.codeHash == codeHash, "INVALID_CODEHASH");
    require(keccak256(result.params) == keccak256(params), "INVALID_PARAMS");

    uint64[] memory decodedInfo = result.data.toUint64List();

    require(uint256(decodedInfo[1]) > lastUpdate, "TIMESTAMP_IS_OLDER_THAN_THE_LAST_UPDATE");

    latestETHPrice = uint256(decodedInfo[0]);
    lastUpdate = uint256(decodedInfo[1]);
  }
}
