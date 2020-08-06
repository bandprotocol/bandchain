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
    pubkeys = []
    len_sigs, remaining = PyObiInteger("u32").decode(signatures)
    for i in range(len_sigs):
        r, remaining = PyObiBytes().decode(remaining)
        s, remaining = PyObiBytes().decode(remaining)
        v, remaining = PyObiInteger("u8").decode(remaining)
        signed_data_prefix, remaining = PyObiBytes().decode(remaining)
        signed_data_suffix, remaining = PyObiBytes().decode(remaining)
        pubkeys.append(
            recover_signer(
                r,
                s,
                v,
                signed_data_prefix,
                signed_data_suffix,
                block_hash
            )
        )
    return pubkeys
