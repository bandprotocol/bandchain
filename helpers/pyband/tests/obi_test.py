import pytest
from pyband.obi import *


def test_encode_decode_bool_success():
    # encode
    assert PyObiBool().encode(True) == bytes.fromhex("01")
    assert PyObiBool().encode(False) == bytes.fromhex("00")

    # decode
    assert PyObiBool().decode(bytes.fromhex("01")) == (True, b"")
    assert PyObiBool().decode(bytes.fromhex("00")) == (False, b"")


def test_encode_decode_bool_fail():
    with pytest.raises(ValueError) as e:
        PyObiBool().decode(bytes.fromhex("07"))

    assert str(e.value) == "Boolean value must be 1 or 0 but got 7"


def test_encode_decode_unsigned_integer_success():
    # encode
    assert PyObiInteger("u8").encode(55) == bytes.fromhex("37")
    assert PyObiInteger("u16").encode(55555) == bytes.fromhex("d903")
    assert PyObiInteger("u32").encode(555555555) == bytes.fromhex("211d1ae3")
    assert PyObiInteger("u64").encode(5555555555555555555) == bytes.fromhex("4d194c57dad638e3")
    assert PyObiInteger("u128").encode(55555555555555555555555555555555555555) == bytes.fromhex(
        "29cb9c5d87a03443cc692f78e38e38e3"
    )
    assert PyObiInteger("u256").encode(
        55555555555555555555555555555555555555555555555555555555555555555555555555555
    ) == bytes.fromhex("7ad35483b719950a2913082cd23202dbb41bd2cb9375038e38e38e38e38e38e3")

    # decode
    assert PyObiInteger("u8").decode(bytes.fromhex("37")) == (55, b"")
    assert PyObiInteger("u16").decode(bytes.fromhex("d903")) == (55555, b"")
    assert PyObiInteger("u32").decode(bytes.fromhex("211d1ae3")) == (555555555, b"")
    assert PyObiInteger("u64").decode(bytes.fromhex("4d194c57dad638e3")) == (
        5555555555555555555,
        b"",
    )
    assert PyObiInteger("u128").decode(bytes.fromhex("29cb9c5d87a03443cc692f78e38e38e3")) == (
        55555555555555555555555555555555555555,
        b"",
    )
    assert PyObiInteger("u256").decode(
        bytes.fromhex("7ad35483b719950a2913082cd23202dbb41bd2cb9375038e38e38e38e38e38e3")
    ) == (55555555555555555555555555555555555555555555555555555555555555555555555555555, b"")


def test_encode_decode_unsigned_integer_fail():
    with pytest.raises(OverflowError) as e:
        PyObiInteger("u8").encode(2 ** 8)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("u16").encode(2 ** 16)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("u32").encode(2 ** 32)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("u64").encode(2 ** 64)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("u128").encode(2 ** 128)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("u256").encode(2 ** 256)

    assert str(e.value) == "int too big to convert"


def test_encode_decode_signed_integer():
    # encode
    assert PyObiInteger("i8").encode(42) == bytes.fromhex("2a")
    assert PyObiInteger("i16").encode(42) == bytes.fromhex("002a")
    assert PyObiInteger("i32").encode(42) == bytes.fromhex("0000002a")
    assert PyObiInteger("i64").encode(42) == bytes.fromhex("000000000000002a")
    assert PyObiInteger("i128").encode(42) == bytes.fromhex("0000000000000000000000000000002a")
    assert PyObiInteger("i256").encode(42) == bytes.fromhex(
        "000000000000000000000000000000000000000000000000000000000000002a"
    )

    assert PyObiInteger("i8").encode(-(2 ** 7)) == bytes.fromhex("80")
    assert PyObiInteger("i16").encode(-(2 ** 15)) == bytes.fromhex("8000")
    assert PyObiInteger("i32").encode(-(2 ** 31)) == bytes.fromhex("80000000")
    assert PyObiInteger("i64").encode(-(2 ** 63)) == bytes.fromhex("8000000000000000")
    assert PyObiInteger("i128").encode(-(2 ** 127)) == bytes.fromhex("80000000000000000000000000000000")
    assert PyObiInteger("i256").encode(-(2 ** 255)) == bytes.fromhex(
        "8000000000000000000000000000000000000000000000000000000000000000"
    )

    # decode
    assert PyObiInteger("i8").decode(bytes.fromhex("2a")) == (42, b"")
    assert PyObiInteger("i16").decode(bytes.fromhex("002a")) == (42, b"")
    assert PyObiInteger("i32").decode(bytes.fromhex("0000002a")) == (42, b"")
    assert PyObiInteger("i64").decode(bytes.fromhex("000000000000002a")) == (42, b"")
    assert PyObiInteger("i128").decode(bytes.fromhex("0000000000000000000000000000002a")) == (
        42,
        b"",
    )
    assert PyObiInteger("i256").decode(
        bytes.fromhex("000000000000000000000000000000000000000000000000000000000000002a")
    ) == (42, b"")

    assert PyObiInteger("i8").decode(bytes.fromhex("80")) == (-(2 ** 7), b"")
    assert PyObiInteger("i16").decode(bytes.fromhex("8000")) == (-(2 ** 15), b"")
    assert PyObiInteger("i32").decode(bytes.fromhex("80000000")) == (-(2 ** 31), b"")
    assert PyObiInteger("i64").decode(bytes.fromhex("8000000000000000")) == (-(2 ** 63), b"")
    assert PyObiInteger("i128").decode(bytes.fromhex("80000000000000000000000000000000")) == (
        -(2 ** 127),
        b"",
    )
    assert PyObiInteger("i256").decode(
        bytes.fromhex("8000000000000000000000000000000000000000000000000000000000000000")
    ) == (-(2 ** 255), b"")


def test_encode_decode_signed_integer_fail():
    with pytest.raises(OverflowError) as e:
        PyObiInteger("i8").encode(2 ** 8)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("i16").encode(2 ** 16)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("i32").encode(2 ** 32)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("i64").encode(2 ** 64)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("i128").encode(2 ** 128)

    assert str(e.value) == "int too big to convert"

    with pytest.raises(OverflowError) as e:
        PyObiInteger("i256").encode(2 ** 256)

    assert str(e.value) == "int too big to convert"


def test_encode_decode_str_success():
    # encode
    assert PyObiString().encode("") == bytes.fromhex("00000000")
    assert PyObiString().encode("mumu") == bytes.fromhex("000000046d756d75")

    # decode
    assert PyObiString().decode(bytes.fromhex("00000000")) == ("", b"")
    assert PyObiString().decode(bytes.fromhex("000000046d756d75")) == ("mumu", b"")


def test_encode_decode_bytes():
    # encode
    assert PyObiBytes().encode(bytes([1, 2, 3, 4, 5, 6, 7, 8, 9])) == bytes.fromhex("00000009010203040506070809")

    # decode
    assert PyObiBytes().decode(bytes.fromhex("00000009010203040506070809")) == (
        bytes([1, 2, 3, 4, 5, 6, 7, 8, 9]),
        b"",
    )


def test_encode_decode_array():
    # encode
    # bool
    assert PyObiArray("[bool;0]").encode([]) == bytes.fromhex("")
    assert PyObiArray("[bool;1]").encode([False]) == bytes.fromhex("00")
    assert PyObiArray("[bool;5]").encode([True, False, True, False, True]) == bytes.fromhex("0100010001")

    # integer
    assert PyObiArray("[i32;0]").encode([]) == bytes.fromhex("")
    assert PyObiArray("[i32;1]").encode([999]) == bytes.fromhex("000003e7")
    assert PyObiArray("[i32;3]").encode([11, 2222, 33333]) == bytes.fromhex("0000000b000008ae00008235")

    # string
    assert PyObiArray("[string;0]").encode([]) == bytes.fromhex("")
    assert PyObiArray("[string;1]").encode(["mumu"]) == bytes.fromhex("000000046d756d75")
    assert PyObiArray("[string;3]").encode(["mumu", "lulu", "mumumumu"]) == bytes.fromhex(
        "000000046d756d75000000046c756c75000000086d756d756d756d75"
    )

    # vector
    assert PyObiArray("[[u64];0]").encode([]) == bytes.fromhex("")
    assert PyObiArray("[[u64];2]").encode([[1, 2, 3, 4, 5], []]) == bytes.fromhex(
        "000000050000000000000001000000000000000200000000000000030000000000000004000000000000000500000000"
    )
    assert PyObiArray("[[string];3]").encode([["a", "", "b"], [], ["aaa", "bbb"]]) == bytes.fromhex(
        "00000003000000016100000000000000016200000000000000020000000361616100000003626262"
    )

    # bytes
    assert PyObiArray("[bytes;0]").encode([]) == bytes.fromhex("")
    assert PyObiArray("[bytes;1]").encode([bytes.fromhex("ababab")]) == bytes.fromhex("00000003ababab")
    assert PyObiArray("[bytes;3]").encode(
        [bytes.fromhex("ababab"), bytes.fromhex("cd"), bytes.fromhex("efef")]
    ) == bytes.fromhex("00000003ababab00000001cd00000002efef")

    # array
    assert PyObiArray("[[bool;2];0]").encode([]) == bytes.fromhex("")
    assert PyObiArray("[[bool;2];3]").encode([[False, False], [True, True], [False, True]]) == bytes.fromhex(
        "000001010001"
    )

    # decode
    # bool
    assert PyObiArray("[bool;0]").decode(bytes.fromhex("")) == ([], b"")
    assert PyObiArray("[bool;1]").decode(bytes.fromhex("00")) == ([False], b"")
    assert PyObiArray("[bool;5]").decode(bytes.fromhex("0100010001")) == (
        [
            True,
            False,
            True,
            False,
            True,
        ],
        b"",
    )

    # integer
    assert PyObiArray("[i32;0]").decode(bytes.fromhex("")) == ([], b"")
    assert PyObiArray("[i32;1]").decode(bytes.fromhex("000003e7")) == ([999], b"")
    assert PyObiArray("[i32;3]").decode(bytes.fromhex("0000000b000008ae00008235")) == (
        [
            11,
            2222,
            33333,
        ],
        b"",
    )

    # string
    assert PyObiArray("[string;0]").decode(bytes.fromhex("")) == ([], b"")
    assert PyObiArray("[string;1]").decode(bytes.fromhex("000000046d756d75")) == (["mumu"], b"")
    assert PyObiArray("[string;3]").decode(
        bytes.fromhex("000000046d756d75000000046c756c75000000086d756d756d756d75")
    ) == (["mumu", "lulu", "mumumumu"], b"")

    # bytes
    assert PyObiArray("[bytes;0]").decode(bytes.fromhex("")) == ([], b"")
    assert PyObiArray("[bytes;1]").decode(bytes.fromhex("00000003ababab")) == (
        [bytes.fromhex("ababab")],
        b"",
    )
    assert PyObiArray("[bytes;3]").decode(bytes.fromhex("00000003ababab00000001cd00000002efef")) == (
        [bytes.fromhex("ababab"), bytes.fromhex("cd"), bytes.fromhex("efef")],
        b"",
    )

    # vector
    assert PyObiArray("[[u64];0]").decode(bytes.fromhex("")) == ([], b"")
    assert PyObiArray("[[u64];2]").decode(
        bytes.fromhex(
            "000000050000000000000001000000000000000200000000000000030000000000000004000000000000000500000000"
        )
    ) == ([[1, 2, 3, 4, 5], []], b"")
    assert PyObiArray("[[string];3]").decode(
        bytes.fromhex("00000003000000016100000000000000016200000000000000020000000361616100000003626262")
    ) == ([["a", "", "b"], [], ["aaa", "bbb"]], b"")

    # array
    assert PyObiArray("[[bool;2];0]").decode(bytes.fromhex("")) == ([], b"")
    assert PyObiArray("[[bool;2];3]").decode(bytes.fromhex("000001010001")) == (
        [
            [False, False],
            [True, True],
            [False, True],
        ],
        b"",
    )


def test_encode_decode_vec():
    # encode
    assert PyObiVector("[i32]").encode([87654321, -12345678]) == bytes.fromhex("0000000205397fb1ff439eb2")
    assert PyObiVector("[string]").encode(["mumu", "imprefvicticious"]) == bytes.fromhex(
        "00000002000000046d756d7500000010696d70726566766963746963696f7573"
    )

    # decode
    assert PyObiVector("[i32]").decode(bytes.fromhex("0000000205397fb1ff439eb2")) == (
        [87654321, -12345678],
        b"",
    )
    assert PyObiVector("[string]").decode(
        bytes.fromhex("00000002000000046d756d7500000010696d70726566766963746963696f7573")
    ) == (["mumu", "imprefvicticious"], b"")


def test_encode_decode_struct():
    # encode
    assert PyObiStruct("""{symbol:string,px:u64}""").encode({"symbol": "BTC", "px": 7777777777777777777}).hex() == (
        "000000034254436bf037ae325f1c71"
    )

    # decode
    assert PyObiStruct("""{symbol:string,px:u64}""").decode(bytes.fromhex("000000034254436bf037ae325f1c71")) == (
        {"symbol": "BTC", "px": 7777777777777777777},
        b"",
    )


def test_decode_multi_obis():
    # BTC, 50, 100
    data = bytes.fromhex("00000003425443000000000000003264")
    symbol, remaining = PyObiString().decode(data)
    x, remaining = PyObiInteger("u64").decode(remaining)
    y, remaining = PyObiInteger("u8").decode(remaining)
    assert ["BTC", 50, 100] == [symbol, x, y]

    # band, 400, 100
    data = bytes.fromhex("0000000462616e64000000000000019064")
    symbol, remaining = PyObiString().decode(data)
    x, remaining = PyObiInteger("u64").decode(remaining)
    y, remaining = PyObiInteger("u8").decode(remaining)
    assert ["band", 400, 100] == [symbol, x, y]


def test_encode_decode_single_success():
    obi = PyObi(
        """
        {
            symbol: string,
            px: u64,
            in: {
                a: u8,
                b: u16,
                in_in: {
                    c:bool,
                    d:bytes
                },
                e: u32
            },
            tb:bool
        }
        """
    )
    test_struct = {
        "symbol": "BTC",
        "px": 9000,
        "in": {
            "a": 1,
            "b": 2,
            "in_in": {"c": True, "d": bytes.fromhex("00112233445566778899")},
            "e": 999999,
        },
        "tb": False,
    }

    encoded = obi.encode(test_struct)
    assert "000000034254430000000000002328010002010000000a00112233445566778899000f423f00" == encoded.hex()
    assert test_struct == obi.decode(encoded)


def test_encode_decode_input_output():
    obi = PyObi("{symbol: string,px: u64,in: {a: u8,b: u8}, tb:bool} / string")
    test_input = {"symbol": "BTC", "px": 9000, "in": {"a": 1, "b": 2}, "tb": False}

    encoded = obi.encode_input(test_input)
    assert "000000034254430000000000002328010200" == encoded.hex()
    assert test_input == obi.decode_input(encoded)

    test_output = "mumumumu"
    encoded = obi.encode_output(test_output)
    assert "000000086d756d756d756d75" == encoded.hex()
    assert test_output == obi.decode_output(encoded)


def test_encode_decode_muti():
    obi = PyObi(
        """
        {
            symbol: string,
            px: u64,
            in: {
                a: u8,
                b: u16,
                in_in: {
                    c:bool,
                    d:bytes
                },
                e: u32
            },
            tb:bool
        }
        /
        string
        /
        u256
        /
        [[u32]]
        /
        {
            a: [
                {
                    a1: [i8],
                    a2: bytes
                }
            ],
            b: {
                b1: [bool]
            },
            c: [string]
        }
        /
        [[u64];10]
        /
        [
            {
                a: [string;2],
                b: u8
            }
            ;3
        ]
        """
    )

    test_structs = [
        {
            "symbol": "BTC",
            "px": 9000,
            "in": {
                "a": 1,
                "b": 2,
                "in_in": {"c": True, "d": bytes.fromhex("00112233445566778899")},
                "e": 999999,
            },
            "tb": False,
        },
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
        11142281140674468977901077842737737931203453356601189394763448166092322549068,
        [[1, 2, 3], [44, 55, 66, 77, 88], [9999, 10000, 10001, 10002, 10003], [12345678], []],
        {
            "a": [
                {"a1": [-1, -2, -3, 0, 1, 2, 3], "a2": bytes.fromhex("abcdef")},
                {"a1": [100], "a2": bytes.fromhex("0" * 10)},
                {"a1": [], "a2": bytes.fromhex("a" * 128)},
            ],
            "b": {"b1": [True, False, True, True, True, True, False]},
            "c": ["mumu", "lulu", "momo", "toto", "coco", "bobo"],
        },
        [[], [1], [2, 22], [3, 33, 333], [4], [5], [6, 66], [7, 7777], [8], [9]],
        [{"a": ["a1", "a2"], "b": 15}, {"a": ["a3", "a4"], "b": 68}, {"a": ["a5", "a6"], "b": 200}],
    ]

    # struct 1
    encoded = obi.encode(test_structs[0], 0)
    assert "000000034254430000000000002328010002010000000a00112233445566778899000f423f00" == encoded.hex()
    assert test_structs[0] == obi.decode(encoded, 0)

    # struct 2
    encoded = obi.encode(test_structs[1], 1)
    assert (
        "000001bd4c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e73656374657475722061646970697363696e6720656c69742c2073656420646f20656975736d6f642074656d706f7220696e6369646964756e74207574206c61626f726520657420646f6c6f7265206d61676e6120616c697175612e20557420656e696d206164206d696e696d2076656e69616d2c2071756973206e6f737472756420657865726369746174696f6e20756c6c616d636f206c61626f726973206e69736920757420616c697175697020657820656120636f6d6d6f646f20636f6e7365717561742e2044756973206175746520697275726520646f6c6f7220696e20726570726568656e646572697420696e20766f6c7570746174652076656c697420657373652063696c6c756d20646f6c6f726520657520667567696174206e756c6c612070617269617475722e204578636570746575722073696e74206f6363616563617420637570696461746174206e6f6e2070726f6964656e742c2073756e7420696e2063756c706120717569206f666669636961206465736572756e74206d6f6c6c697420616e696d20696420657374206c61626f72756d2e"
        == encoded.hex()
    )
    assert test_structs[1] == obi.decode(encoded, 1)

    # struct 3
    encoded = obi.encode(test_structs[2], 2)
    assert "18a24ec1659380f421195719845bbff4df8568c7bd615812b39436bbcb8ead4c" == encoded.hex()
    assert test_structs[2] == obi.decode(encoded, 2)

    # struct 4
    encoded = obi.encode(test_structs[3], 3)
    assert (
        "0000000500000003000000010000000200000003000000050000002c00000037000000420000004d00000058000000050000270f000027100000271100002712000027130000000100bc614e00000000"
        == encoded.hex()
    )
    assert test_structs[3] == obi.decode(encoded, 3)

    # struct 5
    encoded = obi.encode(test_structs[4], 4)
    assert (
        "0000000300000007fffefd0001020300000003abcdef00000001640000000500000000000000000000000040aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa000000070100010101010000000006000000046d756d75000000046c756c75000000046d6f6d6f00000004746f746f00000004636f636f00000004626f626f"
        == encoded.hex()
    )
    assert test_structs[4] == obi.decode(encoded, 4)

    # struct 6
    encoded = obi.encode(test_structs[5], 5)
    assert (
        "0000000000000001000000000000000100000002000000000000000200000000000000160000000300000000000000030000000000000021000000000000014d00000001000000000000000400000001000000000000000500000002000000000000000600000000000000420000000200000000000000070000000000001e61000000010000000000000008000000010000000000000009"
        == encoded.hex()
    )
    assert test_structs[5] == obi.decode(encoded, 5)

    # struct 7
    encoded = obi.encode(test_structs[6], 6)
    print(encoded.hex())
    assert "0000000261310000000261320f00000002613300000002613444000000026135000000026136c8" == encoded.hex()
    assert test_structs[6] == obi.decode(encoded, 6)


def test_encode_decode_not_all_data_is_consumed_fail():
    with pytest.raises(DecodeError) as e:
        PyObi("""bool""").decode(bytes.fromhex("00aa"))

    assert str(e.value) == "Not all data is consumed after decoding input"

    with pytest.raises(DecodeError) as e:
        PyObi("""i64""").decode(bytes.fromhex("00112233445566778899"))

    assert str(e.value) == "Not all data is consumed after decoding input"

    with pytest.raises(DecodeError) as e:
        PyObi("""u64""").decode(bytes.fromhex("00112233445566778899"))

    assert str(e.value) == "Not all data is consumed after decoding input"

    with pytest.raises(DecodeError) as e:
        PyObi("""string""").decode(bytes.fromhex("00000000aabb"))

    assert str(e.value) == "Not all data is consumed after decoding input"

    with pytest.raises(DecodeError) as e:
        PyObi("""bytes""").decode(bytes.fromhex("00000000aabb"))

    assert str(e.value) == "Not all data is consumed after decoding input"

    with pytest.raises(DecodeError) as e:
        PyObi("""[u32]""").decode(bytes.fromhex("00000000aabb"))

    assert str(e.value) == "Not all data is consumed after decoding input"

    with pytest.raises(DecodeError) as e:
        PyObi("""{ x:i16 }""").decode(bytes.fromhex("00000000"))

    assert str(e.value) == "Not all data is consumed after decoding input"
