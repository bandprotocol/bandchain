pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

interface IBridge {
	/// Helper struct to help the function caller to decode oracle data.
	struct VerifyOracleDataResult {
		bytes data;
		bytes32 codeHash;
		bytes params;
	}

	function relayAndVerify(bytes calldata _data)
		external
		returns (VerifyOracleDataResult memory result);
}
