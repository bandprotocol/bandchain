pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {SafeMath} from "openzeppelin-solidity/contracts/math/SafeMath.sol";
import {Ownable} from "openzeppelin-solidity/contracts/ownership/Ownable.sol";
import {IBridge} from "../IBridge.sol";

contract DEX is Ownable {
    using SafeMath for uint256;

    uint64 public oracleScriptId;

    mapping(address => mapping(bytes => uint256)) _balances;
    mapping(string => bytes) public supportedTokens;

    IBridge bridge;

    constructor(IBridge _bridge, uint64 _oracleScriptId) public {
        bridge = _bridge;
        oracleScriptId = _oracleScriptId;
        supportedTokens["ADA"] = hex"00000000";
        supportedTokens["BAND"] = hex"00000001";
        supportedTokens["BCH"] = hex"00000002";
        supportedTokens["BNB"] = hex"00000003";
        supportedTokens["BTC"] = hex"00000004";
        supportedTokens["EOS"] = hex"00000005";
        supportedTokens["ETC"] = hex"00000006";
        supportedTokens["ETH"] = hex"00000007";
        supportedTokens["LTC"] = hex"00000008";
        supportedTokens["TRX"] = hex"00000009";
        supportedTokens["XRP"] = hex"0000000A";
    }

    function isSupportedToken(string memory _a)
        public
        view
        returns (bytes memory, bool)
    {
        bytes memory key = supportedTokens[_a];
        return (key, key.length != 0);
    }

    function balanceOf(address _account, string memory _symbol)
        public
        view
        returns (uint256)
    {
        (bytes memory key, bool _isSupportedToken) = isSupportedToken(_symbol);
        if (!_isSupportedToken) {
            revert("UNKNOWN_SYMBOL");
        }
        return _balances[_account][key];
    }

    function bytesToPrices(bytes memory _b)
        public
        pure
        returns (uint256, uint256)
    {
        require(_b.length == 16, "INVALID_LENGTH");
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
        IBridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
            _reportPrice
        );

        require(
            result.oracleScriptId == oracleScriptId,
            "INVALID_ORACLE_SCRIPT"
        );

        (uint256 ethPrice, uint256 otherPrice) = bytesToPrices(result.data);

        uint256 tokenEarn = msg.value.mul(ethPrice).div(otherPrice);

        _balances[msg.sender][result.params] = _balances[msg.sender][result
            .params]
            .add(tokenEarn);
    }

    function sell(uint256 _amount, bytes memory _reportPrice) public {
        IBridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
            _reportPrice
        );

        require(
            result.oracleScriptId == oracleScriptId,
            "INVALID_ORACLE_SCRIPT"
        );
        require(
            _amount <= _balances[msg.sender][result.params],
            "INSUFFICIENT_TOKENS"
        );

        (uint256 ethPrice, uint256 otherPrice) = bytesToPrices(result.data);

        uint256 ethEarn = _amount.mul(otherPrice).div(ethPrice);

        _balances[msg.sender][result.params] = _balances[msg.sender][result
            .params]
            .sub(_amount);
        msg.sender.transfer(ethEarn);
    }

    function withdraw() public onlyOwner {
        msg.sender.transfer(address(this).balance);
    }
}
