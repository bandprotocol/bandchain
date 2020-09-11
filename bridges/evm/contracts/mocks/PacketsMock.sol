// SPDX-License-Identifier: Apache-2.0

pragma solidity 0.6.11;
pragma experimental ABIEncoderV2;

import {IBridge} from "../interfaces/IBridge.sol";
import {Packets} from "../Packets.sol";


contract PacketsMock {
    function encodeRequestPacket(IBridge.RequestPacket memory _packet)
        public
        pure
        returns (bytes memory)
    {
        return Packets.encodeRequestPacket(_packet);
    }

    function encodeResponsePacket(IBridge.ResponsePacket memory _packet)
        public
        pure
        returns (bytes memory)
    {
        return Packets.encodeResponsePacket(_packet);
    }

    function getEncodedResult(
        IBridge.RequestPacket memory _req,
        IBridge.ResponsePacket memory _res
    ) public pure returns (bytes memory) {
        return Packets.getEncodedResult(_req, _res);
    }
}
