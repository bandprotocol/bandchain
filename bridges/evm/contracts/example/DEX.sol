pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import { SafeMath } from "openzeppelin-solidity/contracts/math/SafeMath.sol";
import { IBridge } from "../IBridge.sol";

contract DEX {
  using SafeMath for uint256;

  bytes32 public codeHash;

  mapping(address => mapping(bytes => uint256)) _balances;

  IBridge bridge;

  constructor(IBridge _bridge, bytes32 _codeHash) public {
    bridge = _bridge;
    codeHash = _codeHash;
  }

  function strCmp(string memory a, string memory b)
    public
    pure
    returns (bool)
  {
    return keccak256(abi.encodePacked(a)) == keccak256(abi.encodePacked(b));
  }

  function balanceOf(address account, string memory symbol)
    public
    view
    returns (uint256)
  {
    bytes memory key;
    if (strCmp(symbol, "ADA")) {
      key = hex"00000000";
    } else if (strCmp(symbol, "BAND")) {
      key = hex"00000001";
    } else if (strCmp(symbol, "BCH")) {
      key = hex"00000002";
    } else if (strCmp(symbol, "BNB")) {
      key = hex"00000003";
    } else if (strCmp(symbol, "BTC")) {
      key = hex"00000004";
    } else if (strCmp(symbol, "EOS")) {
      key = hex"00000005";
    } else if (strCmp(symbol, "ETC")) {
      key = hex"00000006";
    } else if (strCmp(symbol, "ETH")) {
      key = hex"00000007";
    } else if (strCmp(symbol, "LTC")) {
      key = hex"00000008";
    } else if (strCmp(symbol, "TRX")) {
      key = hex"00000009";
    } else if (strCmp(symbol, "XRP")) {
      key = hex"0000000A";
    } else {
      revert("UNKNOWN_SYMBOL");
    }
    return _balances[account][key];
  }

  function bytesToPrices(bytes memory _b)
    public
    pure
    returns (uint256, uint256)
  {
    require(_b.length >= 16, "INVALID_LENGTH");
    uint256 ethPrice;
    uint256 otherPrice;
    for (uint256 i = 0; i < 8; i++) {
      ethPrice = ethPrice + (uint256(uint8(_b[i])) << (8 * (7 - i)));
      otherPrice = otherPrice + (uint256(uint8(_b[i + 8])) << (8 * (7 - i)));
    }
    return (ethPrice, otherPrice);
  }

  function buy(bytes memory _reportPrice) public payable {
    IBridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
      _reportPrice
    );

    require(result.codeHash == codeHash, "INVALID_CODEHASH");

    (uint256 ethPrice, uint256 otherPrice) = bytesToPrices(result.data);

    uint256 tokenEarn = msg.value.mul(ethPrice).div(otherPrice);

    _balances[msg.sender][result.params] = _balances[msg.sender][result.params].add(tokenEarn);
  }

  function sell(uint256 amount, bytes memory _reportPrice) public {
    IBridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
      _reportPrice
    );

    require(result.codeHash == codeHash, "INVALID_CODEHASH");
    require(
      amount <= _balances[msg.sender][result.params],
      "INSUFFICIENT_TOKENS"
    );

    (uint256 ethPrice, uint256 otherPrice) = bytesToPrices(result.data);

    uint256 ethEarn = amount.mul(otherPrice).div(ethPrice);

    _balances[msg.sender][result.params] = _balances[msg.sender][result.params].sub(amount);
    msg.sender.transfer(ethEarn);
  }

  function withdraw() public {
    msg.sender.transfer(address(this).balance);
  }
}
