pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

interface Bridge {
  struct VerifyOracleDataResult {
    bytes data;
    bytes32 codeHash;
    bytes params;
  }
  function relayAndVerify(bytes calldata data) external returns (VerifyOracleDataResult memory result);
}

contract NeoTCD {
  bytes32 public codeHash;
  bytes public params;
  uint256 public latestETHPrice;
  uint256 public latestReportedBlock;

  Bridge public bridge;

  constructor(bytes32 _codeHash , bytes memory _params, address bridgeAddress) public {
    codeHash = _codeHash;
    params = _params;
    bridge = Bridge(bridgeAddress);
  }

  function update(bytes memory _reportPrice) public {
    Bridge.VerifyOracleDataResult memory result = bridge.relayAndVerify(_reportPrice);

    require(result.codeHash == codeHash, "INVALID_CODEHASH");
    require(keccak256(result.params) == keccak256(params), "INVALID_PARAMS");

    uint256 _latestETHPrice;
    uint256 timestamp;
    uint256 halfSize = result.data.length/2;
    for(uint256 i = 0; i < halfSize; i++){
      uint256 j = i + halfSize;
      uint256 mantissa = (2**(8*(halfSize-(i+1))));
      _latestETHPrice = _latestETHPrice + uint256(uint8(result.data[i]))*mantissa;
      timestamp = timestamp + uint256(uint8(result.data[j]))*mantissa;
    }

    uint256 dt;
    if (timestamp < now) {
      dt = now - timestamp;
    } else {
      dt = timestamp - now;
    }

    require(dt < 3 minutes);

    latestETHPrice = _latestETHPrice;
    latestReportedBlock = timestamp;
  }
}
