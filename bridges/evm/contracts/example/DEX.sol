pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {SafeMath} from "openzeppelin-solidity/contracts/math/SafeMath.sol";

interface Bridge {
    struct VerifyOracleDataResult {
        bytes data;
        bytes32 codeHash;
        bytes params;
    }
    function relayAndVerify(bytes calldata data)
        external
        returns (VerifyOracleDataResult memory result);
}

contract DEX {
    using SafeMath for uint256;

    bytes32 public codeHash;

    mapping(address => mapping(bytes => uint256)) public balances;

    Bridge bridge = Bridge(0x3e1F8745E4088443350121075828F119075ef641);

    constructor(bytes32 _codeHash) public {
        codeHash = _codeHash;
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
            otherPrice =
                otherPrice +
                (uint256(uint8(_b[i + 8])) << (8 * (7 - i)));
        }
        return (ethPrice, otherPrice);
    }

    function buy(bytes memory _reportPrice) public payable {
        Bridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
            _reportPrice
        );

        require(result.codeHash == codeHash, "INVALID_CODEHASH");

        (uint256 ethPrice, uint256 otherPrice) = bytesToPrices(result.data);

        uint256 tokenEarn = msg.value.mul(ethPrice).div(otherPrice);

        balances[msg.sender][result.params] = balances[msg.sender][result
            .params]
            .add(tokenEarn);
    }

    function sell(uint256 amount, bytes memory _reportPrice) public {
        Bridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
            _reportPrice
        );

        require(result.codeHash == codeHash, "INVALID_CODEHASH");
        require(
            amount <= balances[msg.sender][result.params],
            "INSUFFICIENT_TOKENS"
        );

        (uint256 ethPrice, uint256 otherPrice) = bytesToPrices(result.data);

        uint256 ethEarn = amount.mul(otherPrice).div(ethPrice);

        balances[msg.sender][result.params] = balances[msg.sender][result
            .params]
            .sub(amount);
        msg.sender.transfer(ethEarn);
    }

    function withdraw() public {
        msg.sender.transfer(address(this).balance);
    }
}
