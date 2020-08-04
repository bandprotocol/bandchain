# OBI minimal
from iconservice import *


def decode_int(b: bytes, n_bits: int) -> (int, bytes):
    if n_bits <= 0:
        revert("INVALID_INT_FORMAT")

    n = n_bits // 8
    if n * 8 != n_bits:
        revert("INVALID_INT_FORMAT")

    if n > len(b):
        revert("DECODE_INT_SIZE_EXCEED")

    return (int.from_bytes(b[:n], 'big'), b[n:])


def decode_bool(b: bytes) -> (int, bytes):
    val, remaining = decode_int(b, 8)
    if val == 0:
        return (False, remaining)
    elif val == 1:
        return (True, remaining)

    revert("INVALID_BOOL_FORMAT")


def decode_bytes(b: bytes) -> (bytes, bytes):
    (size, remaining) = decode_int(b, 32)
    return (remaining[:size], remaining[size:])


def decode_str(b: bytes) -> (str, bytes):
    (size, remaining) = decode_int(b, 32)
    return (remaining[:size].decode("utf-8"), remaining[size:])
