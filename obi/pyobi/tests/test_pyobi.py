import pytest
from pyobi import *


def test_encode_signed_integer():
    assert PyObiInteger("i8").encode(42) == bytes.fromhex("2a")
    assert PyObiInteger("i16").encode(42) == bytes.fromhex("002a")
    assert PyObiInteger("i32").encode(42) == bytes.fromhex("0000002a")
    assert PyObiInteger("i64").encode(42) == bytes.fromhex("000000000000002a")
    assert PyObiInteger("i128").encode(42) == bytes.fromhex("0000000000000000000000000000002a")
    assert PyObiInteger("i256").encode(42) == bytes.fromhex(
        "000000000000000000000000000000000000000000000000000000000000002a"
    )
