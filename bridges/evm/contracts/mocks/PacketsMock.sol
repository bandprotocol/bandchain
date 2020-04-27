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
}
