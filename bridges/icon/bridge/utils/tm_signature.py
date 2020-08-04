from iconservice import *
from . import secp256k1, sha256, obi


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
    len_sigs, remaining = obi.decode_int(signatures, 32)
    for i in range(len_sigs):
        r, remaining = obi.decode_bytes(remaining)
        s, remaining = obi.decode_bytes(remaining)
        v, remaining = obi.decode_int(remaining, 8)
        signed_data_prefix, remaining = obi.decode_bytes(remaining)
        signed_data_suffix, remaining = obi.decode_bytes(remaining)
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
