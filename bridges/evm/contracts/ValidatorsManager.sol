pragma solidity 0.5.14;
import { Ownable } from "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract ValidatorsManager is Ownable {
  uint256 public numberOfValidators;
  mapping(address => bool) public validators;

  constructor(address[] memory _initValidators) public {
    addValidators(_initValidators);
  }

  function verifyValidators(address[] memory addrs) public view returns(bool) {
    require(addrs.length * 3 > numberOfValidators * 2);
    for (uint256 i = 0; i < addrs.length; i++) {
      require(validators[addrs[i]], "INVALID_VALIDATOR_ADDRESS");
    }
    return true;
  }

  function addValidators(address[] memory _validators) public onlyOwner {
    for (uint i = 0; i < _validators.length; i++) {
      require(!validators[_validators[i]], "ALREADY_BE_VALIDATOR");
      validators[_validators[i]] = true;
    }
    numberOfValidators += _validators.length;
  }

  function removeValidators(address[] memory _validators) public onlyOwner {
    for (uint i = 0; i < _validators.length; i++) {
      require(validators[_validators[i]], "UNKNOWN_VALIDATOR");
      validators[_validators[i]] = false;
    }
    numberOfValidators -= _validators.length;
  }
}
