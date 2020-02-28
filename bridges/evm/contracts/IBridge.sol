pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

interface IBridge {
  /// Helper struct to help wraping the information of result and its context.
  struct WrappedResult {
    uint64 requestTime;
    uint64 aggregationTime;
    uint64 requestedValidatorsCount;
    uint64 sufficientValidatorCount;
    uint64 reportedValidatorsCount;
    bytes data;
  }

  /// Helper struct to help the function caller to decode oracle data.
  struct VerifyOracleDataResult {
    WrappedResult result;
    uint64 oracleScriptId;
    bytes params;
  }

  /// Performs oracle state relay and oracle data verification in one go. The caller submits
  /// the encoded proof and receives back the decoded data, ready to be validated and used.
  /// @param _data The encoded data for oracle state relay and data verification.
  function relayAndVerify(bytes calldata _data)
    external
    returns (VerifyOracleDataResult memory result);
}
