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
                uint32(bytes(_self.clientId).length),
                _self.clientId,
                _self.oracleScriptId,
                uint32(_self.params.length),
                _self.params,
                _self.askCount,
                _self.minCount
            );
    }

    function marshalResponsePacket(IBridge.ResponsePacket memory _self)
        internal
        pure
        returns (bytes memory)
    {
        return
            abi.encodePacked(
                uint32(bytes(_self.clientId).length),
                _self.clientId,
                _self.requestId,
                _self.ansCount,
                _self.requestTime,
                _self.resolveTime,
                uint32(_self.resolveStatus),
                uint32(bytes(_self.result).length),
                _self.result
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
