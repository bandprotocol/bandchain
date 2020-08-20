from iconservice import *
from . import sha256


def merkle_leaf_hash(value: bytes) -> bytes:
    return sha256.digest(bytes([0]) + value)


def merkle_inner_hash(left: bytes, right: bytes) -> bytes:
    return sha256.digest(bytes([1]) + left + right)


def encode_varint_unsigned(value: int) -> bytes:
    result = b""
    while value > 0:
        result += bytes([128 | (value & 127)])
        value >>= 7

    return result[: len(result) - 1] + bytes([result[len(result) - 1] & 127])


def encode_varint_signed(value: int) -> bytes:
    return encode_varint_unsigned(value * 2)
