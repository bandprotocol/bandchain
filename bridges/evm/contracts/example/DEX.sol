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

contract TicketSeller {
    using SafeMath for uint256;

    bytes32 public codeHash;
    bytes public params;

    uint256 public priceInUSD;

    uint64 public currentEthPrice;

    mapping(address => uint256) public tickets;

    Bridge bridge = Bridge(0x3e1F8745E4088443350121075828F119075ef641);

    constructor(bytes32 _codeHash, bytes memory _params, uint256 _price)
        public
    {
        codeHash = _codeHash;
        params = _params;
        priceInUSD = _price;
    }

    function bytesToUInt(bytes memory _b) public pure returns (uint256) {
        uint256 number;
        for (uint256 i = 0; i < _b.length; i++) {
            number =
                number +
                uint256(uint8(_b[i])) *
                (2**(8 * (_b.length - (i + 1))));
        }
        return number;
    }

    function buy(bytes memory _reportPrice) public payable {
        Bridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(
            _reportPrice
        );

        require(result.codeHash == codeHash, "INVALID_CODEHASH");
        require(
            keccak256(result.params) == keccak256(params),
            "INVALID_PARAMS"
        );

        uint256 ethPrice = bytesToUInt(result.data);

        uint256 ticketPrice = priceInUSD.mul(1e20).div(ethPrice);

        require(msg.value >= ticketPrice);

        if (msg.value > ticketPrice) {
            msg.sender.transfer(msg.value - ticketPrice);
        }
        tickets[msg.sender]++;
    }

    function withdraw() public {
        msg.sender.transfer(address(this).balance);
    }
}
