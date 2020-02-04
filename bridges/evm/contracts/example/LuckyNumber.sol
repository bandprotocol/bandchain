pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import { Ownable } from "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import { IBridge } from "bandchain-helper-library/contracts/IBridge.sol";

contract LuckyNumber is Ownable {
  bytes32 public codeHash;
  uint64 public maxNumber;

  uint64 target;

  IBridge bridge;

  constructor(IBridge _bridge, bytes32 _codeHash, uint64 _maxNumber, uint64 _target)
    public
    payable
  {
    bridge = _bridge;
    codeHash = _codeHash;
    maxNumber = _maxNumber;
    target = _target;
  }

  function bytesToU64(bytes memory _b) public pure returns (uint64) {
    require(_b.length == 8, "INVALID_LENGTH");
    uint64 number;
    for (uint256 i = 0; i < 8; i++) {
      number = number + (uint64(uint8(_b[i])) << (8 * (7 - i)));
    }
    return number;
  }

  function guess(bytes memory _reportPrice) public {
    IBridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
      _reportPrice
    );

    require(result.codeHash == codeHash, "INVALID_CODEHASH");

    require(maxNumber == bytesToU64(result.params), "INVALID_MAX_NUM");
    require(target == bytesToU64(result.data), "WRONG_GUESS");
    msg.sender.transfer(address(this).balance);
  }

  function withdraw() public onlyOwner {
    msg.sender.transfer(address(this).balance);
  }

  function () external payable {}
}
