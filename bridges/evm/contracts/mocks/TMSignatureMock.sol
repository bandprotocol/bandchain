pragma solidity 0.5.14;
pragma experimental ABIEncoderV2;

import {TMSignature} from "../TMSignature.sol";


contract TMSignatureMock {
    function recoverSigner(
        TMSignature.Data memory _data,
        bytes32 _blockHash,
        bytes memory _signedDataPrefix
    ) public pure returns (address) {
        return TMSignature.recoverSigner(_data, _blockHash, _signedDataPrefix);
    }
}
