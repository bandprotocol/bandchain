from iconservice import *
from . import secp256k1, sha256
from ..pyobi import *


def recover_signer(
    r: bytes,
    s: bytes,
    v: int,
    signed_data_prefix: bytes,
    signed_data_suffix: bytes,
    block_hash: bytes
) -> bytes:
    return secp256k1.ecrecover(sha256.digest(signed_data_prefix+block_hash+signed_data_suffix), r, s, v)


def recover_signers(signatures: bytes, block_hash: bytes) -> list:
    obi = PyObi(
        """
        [
            {
                r: bytes,
                s: bytes,
                v: u8,
                signed_data_prefix: bytes,
                signed_data_suffix: bytes
            }
        ]
        """
    )

    return [
        recover_signer(
            sig["r"],
            sig["s"],
            sig["v"],
            sig["signed_data_prefix"],
            sig["signed_data_suffix"],
            block_hash
        ) for sig in obi.decode(signatures)
    ]
