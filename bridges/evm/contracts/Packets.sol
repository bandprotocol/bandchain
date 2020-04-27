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
                hex"d9c58927",
                uint8((1 << 3) | 2),
                Utils.encodeVarintUnsigned(bytes(_self.clientId).length),
                _self.clientId,
                uint8((2 << 3) | 0),
                Utils.encodeVarintUnsigned(_self.oracleScriptId),
                uint8((3 << 3) | 2),
                Utils.encodeVarintUnsigned(bytes(_self.params).length),
                _self.params,
                uint8((4 << 3) | 0),
                Utils.encodeVarintUnsigned(_self.askCount),
                uint8((5 << 3) | 0),
                Utils.encodeVarintUnsigned(_self.minCount)
            );
    }
}
