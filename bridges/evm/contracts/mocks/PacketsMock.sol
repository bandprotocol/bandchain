pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {IBridge} from "../IBridge.sol";
import {Packets} from "../Packets.sol";


contract PacketsMock {
    function marshalRequestPacket(IBridge.RequestPacket memory _packet)
        public
        pure
        returns (bytes memory)
    {
        return Packets.marshalRequestPacket(_packet);
    }

    function marshalResponsePacket(IBridge.ResponsePacket memory _packet)
        public
        pure
        returns (bytes memory)
    {
        return Packets.marshalResponsePacket(_packet);
    }

    function getResultHash(
        IBridge.RequestPacket memory _req,
        IBridge.ResponsePacket memory _res
    ) public pure returns (bytes32) {
        return Packets.getResultHash(_req, _res);
    }
}
