from iconservice import *
from . import sha256


def merkle_leaf_hash(value: bytes) -> bytes:
    return sha256.digest(bytes([0]) + value)


def merkle_inner_hash(left: bytes, right: bytes) -> bytes:
    return sha256.digest(bytes([1]) + left + right)


def encode_varint_unsigned(value: int) -> bytes:
    temp_value = value
    size = 0
    while temp_value > 0:
        size += 1
        temp_value >>= 7

    result = b''
    temp_value = value
    for i in range(size):
        result += bytes([128 | (temp_value & 127)])
        temp_value >>= 7

    return result[:size - 1] + bytes([result[size - 1] & 127])


def encode_varint_signed(value: int) -> bytes:
    return encode_varint_unsigned(value * 2)
