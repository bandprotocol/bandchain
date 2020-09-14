// SPDX-License-Identifier: Apache-2.0
pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

abstract contract IStdReference {
    /// A structure returned whenever someone requests for standard reference data.
    struct ReferenceData {
        uint256 rate; // base/quote exchange rate, multiplied by 1e18.
        uint256 lastUpdatedBase; // UNIX epoch of the last time when base price gets updated.
        uint256 lastUpdatedQuote; // UNIX epoch of the last time when quote price gets updated.
    }

    /// Returns the price data for the given base/quote pair. Revert if not available.
    function getReferenceData(string memory _base, string memory _quote)
        public
        virtual
        view
        returns (ReferenceData memory);

    /// Similar to getReferenceData, but with multiple base/quote pairs at once.
    function getRefenceDataBulk(string[] memory _bases, string[] memory _quotes)
        public
        view
        returns (ReferenceData[] memory)
    {
        require(_bases.length == _quotes.length, "BAD_INPUT_LENGTH");
        uint256 len = _bases.length;
        ReferenceData[] memory results = new ReferenceData[](len);
        for (uint256 idx = 0; idx < len; idx++) {
            results[idx] = getReferenceData(_bases[idx], _quotes[idx]);
        }
        return results;
    }
}
