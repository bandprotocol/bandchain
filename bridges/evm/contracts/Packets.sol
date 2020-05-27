pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {Utils} from "./Utils.sol";
import {IBridge} from "./IBridge.sol";


library Packets {
    function marshalRequestPacket(IBridge.RequestPacket memory _self)
        internal
        pure
        returns (bytes memory)
    {
        return
            abi.encodePacked(
                bytes4(uint32(bytes(_self.clientId).length)),
                _self.clientId,
                bytes8(_self.oracleScriptId),
                bytes4(uint32(_self.params.length)),
                _self.params,
                bytes8(_self.askCount),
                bytes8(_self.minCount)
            );
    }

    function getEncodeLength(string memory _s)
        internal
        pure
        returns (bytes memory)
    {
        return Utils.encodeVarintUnsigned(bytes(_s).length);
    }

    function getReponsePart1(
        string memory _clientId,
        uint64 _requestId,
        uint64 _ansCount,
        uint64 _requestTime
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(
                uint8(10), // (1 << 3) | 2
                getEncodeLength(_clientId),
                _clientId,
                uint8(16), // (2 << 3) | 0
                Utils.encodeVarintUnsigned(_requestId),
                uint8(24), // (3 << 3) | 0
                Utils.encodeVarintUnsigned(_ansCount),
                uint8(32), // (4 << 3) | 0
                Utils.encodeVarintUnsigned(_requestTime)
            );
    }

    function getReponsePart2(
        uint64 _resolveTime,
        uint8 _resolveStatus,
        bytes memory _result
    ) internal pure returns (bytes memory) {
        return
            abi.encodePacked(
                uint8(40), // (5 << 3) | 0
                Utils.encodeVarintUnsigned(_resolveTime),
                uint8(48), // (6 << 3) | 0
                Utils.encodeVarintUnsigned(_resolveStatus),
                uint8(58), // (7 << 3) | 2
                Utils.encodeVarintUnsigned(_result.length),
                _result
            );
    }

    function marshalResponsePacket(IBridge.ResponsePacket memory _self)
        internal
        pure
        returns (bytes memory)
    {
        return
            abi.encodePacked(
                hex"79b5957c", // Amino codec info for ResponsePacket type
                getReponsePart1(
                    _self.clientId,
                    _self.requestId,
                    _self.ansCount,
                    _self.requestTime
                ),
                getReponsePart2(
                    _self.resolveTime,
                    _self.resolveStatus,
                    _self.result
                )
            );
    }

    function getResultHash(
        IBridge.RequestPacket memory _req,
        IBridge.ResponsePacket memory _res
    ) internal pure returns (bytes32) {
        return
            sha256(
                abi.encodePacked(
                    sha256(marshalRequestPacket(_req)),
                    sha256(marshalResponsePacket(_res))
                )
            );
    }
}
