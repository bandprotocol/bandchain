pragma solidity 0.5.14;

interface IPriceReference {
    /// @dev Returns the number of times that the price has been updated.
    function latestRound() external view returns (uint256);

    /// @dev Returns the latest available price data point.
    function latestAnswer() external view returns (uint256);

    /// @dev Returns the timestamp associated with the latest available price.
    function latestTimestamp() external view returns (uint256);

    /// @dev Returns the price data point from the given round.
    function getAnswer(uint256 _round) external view returns (uint256);

    /// @dev Returns the timestamp associated with the price at the given round.
    function getTimestamp(uint256 _round) external view returns (uint256);
}
