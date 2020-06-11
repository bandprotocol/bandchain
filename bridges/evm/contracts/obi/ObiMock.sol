pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

import {Obi} from "./Obi.sol";
import {ResultDecoder} from "./Result.sol";


contract ObiMock {
    using ResultDecoder for bytes;

    function decodeU8(Obi.Data memory _data) public pure returns (uint8) {
        return Obi.decodeU8(_data);
    }

    function decodeI8(Obi.Data memory _data) public pure returns (int8) {
        return Obi.decodeI8(_data);
    }

    function decodeU16(Obi.Data memory _data) public pure returns (uint16) {
        return Obi.decodeU16(_data);
    }

    function decodeI16(Obi.Data memory _data) public pure returns (int16) {
        return Obi.decodeI16(_data);
    }

    function decodeU32(Obi.Data memory _data) public pure returns (uint32) {
        return Obi.decodeU32(_data);
    }

    function decodeI32(Obi.Data memory _data) public pure returns (int32) {
        return Obi.decodeI32(_data);
    }

    function decodeU64(Obi.Data memory _data) public pure returns (uint64) {
        return Obi.decodeU64(_data);
    }

    function decodeI64(Obi.Data memory _data) public pure returns (int64) {
        return Obi.decodeI64(_data);
    }

    function decodeU128(Obi.Data memory _data) public pure returns (uint128) {
        return Obi.decodeU128(_data);
    }

    function decodeI128(Obi.Data memory _data) public pure returns (int128) {
        return Obi.decodeI128(_data);
    }

    function decodeU256(Obi.Data memory _data) public pure returns (uint256) {
        return Obi.decodeU256(_data);
    }

    function decodeI256(Obi.Data memory _data) public pure returns (int256) {
        return Obi.decodeI256(_data);
    }

    function decodeBool(Obi.Data memory _data) public pure returns (bool) {
        return Obi.decodeBool(_data);
    }

    function decodeBytes(Obi.Data memory _data)
        public
        pure
        returns (bytes memory)
    {
        return Obi.decodeBytes(_data);
    }
}
