from ..pyobi import *
from tbears.libs.scoretest.score_test_case import ScoreTestCase
from iconservice.base.exception import IconScoreException


class TestOBI(ScoreTestCase):

    def setUp(self):
        super().setUp()

    def test_encode_signed_integer(self):
        self.assertEqual(PyObiInteger("i8").encode(42), bytes.fromhex("2a"))
        self.assertEqual(PyObiInteger("i16").encode(42), bytes.fromhex("002a"))
        self.assertEqual(
            PyObiInteger("i32").encode(42),
            bytes.fromhex("0000002a")
        )
        self.assertEqual(
            PyObiInteger("i64").encode(42),
            bytes.fromhex("000000000000002a")
        )
        self.assertEqual(
            PyObiInteger("i128").encode(42),
            bytes.fromhex("0000000000000000000000000000002a")
        )
        self.assertEqual(
            PyObiInteger("i256").encode(42),
            bytes.fromhex(
                "000000000000000000000000000000000000000000000000000000000000002a"
            )
        )

    def test_chain_encoding(self):
        # BTC, 50, 100
        data = bytes.fromhex("00000003425443000000000000003264")
        symbol, remaining = PyObiString().decode(data)
        x, remaining = PyObiInteger("u64").decode(remaining)
        y, remaining = PyObiInteger("u8").decode(remaining)
        self.assertEqual(["BTC", 50, 100], [symbol, x, y])

        # band, 400, 100
        data = bytes.fromhex("0000000462616e64000000000000019064")
        symbol, remaining = PyObiString().decode(data)
        x, remaining = PyObiInteger("u64").decode(remaining)
        y, remaining = PyObiInteger("u8").decode(remaining)
        self.assertEqual(["band", 400, 100], [symbol, x, y])

    def test_encode_decode_a_struct(self):
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
                "in_in": {
                    "c": True,
                    "d": bytes.fromhex("00112233445566778899")
                },
                "e": 999999,
            },
            "tb": False
        }

        encoded = obi.encode(test_struct)
        self.assertEqual(
            "000000034254430000000000002328010002010000000a00112233445566778899000f423f00", encoded.hex())
        self.assertEqual(test_struct, obi.decode(encoded))

    def test_encode_decode_input_output(self):
        obi = PyObi(
            "{symbol: string,px: u64,in: {a: u8,b: u8}, tb:bool} / string")
        test_input = {
            "symbol": "BTC",
            "px": 9000,
            "in": {
                "a": 1,
                "b": 2,
            },
            "tb": False
        }

        encoded = obi.encode_input(test_input)
        self.assertEqual("000000034254430000000000002328010200", encoded.hex())
        self.assertEqual(test_input, obi.decode_input(encoded))

        test_output = "mumumumu"
        encoded = obi.encode_output(test_output)
        self.assertEqual("000000086d756d756d756d75", encoded.hex())
        self.assertEqual(test_output, obi.decode_output(encoded))

    def test_encode_decode_muti_structs(self):
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
            """
        )

        test_structs = [
            {
                "symbol": "BTC",
                "px": 9000,
                "in": {
                    "a": 1,
                    "b": 2,
                    "in_in": {
                        "c": True,
                        "d": bytes.fromhex("00112233445566778899")
                    },
                    "e": 999999,
                },
                "tb": False
            },
            "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
            11142281140674468977901077842737737931203453356601189394763448166092322549068,
            [
                [1, 2, 3],
                [44, 55, 66, 77, 88],
                [9999, 10000, 10001, 10002, 10003],
                [12345678],
                []
            ],
            {
                "a": [
                    {
                        "a1": [-1, -2, -3, 0, 1, 2, 3],
                        "a2": bytes.fromhex("abcdef")
                    },
                    {
                        "a1": [100],
                        "a2": bytes.fromhex("0"*10)
                    },
                    {
                        "a1": [],
                        "a2": bytes.fromhex("a"*128)
                    }
                ],
                "b": {
                    "b1": [True, False, True, True, True, True, False]
                },
                "c": ["mumu", "lulu", "momo", "toto", "coco", "bobo"]
            }
        ]

        # sturct 1
        encoded = obi.encode(test_structs[0], 0)
        self.assertEqual(
            "000000034254430000000000002328010002010000000a00112233445566778899000f423f00", encoded.hex())
        self.assertEqual(test_structs[0], obi.decode(encoded, 0))

        # sturct 2
        encoded = obi.encode(test_structs[1], 1)
        self.assertEqual(
            "000001bd4c6f72656d20697073756d20646f6c6f722073697420616d65742c20636f6e73656374657475722061646970697363696e6720656c69742c2073656420646f20656975736d6f642074656d706f7220696e6369646964756e74207574206c61626f726520657420646f6c6f7265206d61676e6120616c697175612e20557420656e696d206164206d696e696d2076656e69616d2c2071756973206e6f737472756420657865726369746174696f6e20756c6c616d636f206c61626f726973206e69736920757420616c697175697020657820656120636f6d6d6f646f20636f6e7365717561742e2044756973206175746520697275726520646f6c6f7220696e20726570726568656e646572697420696e20766f6c7570746174652076656c697420657373652063696c6c756d20646f6c6f726520657520667567696174206e756c6c612070617269617475722e204578636570746575722073696e74206f6363616563617420637570696461746174206e6f6e2070726f6964656e742c2073756e7420696e2063756c706120717569206f666669636961206465736572756e74206d6f6c6c697420616e696d20696420657374206c61626f72756d2e", encoded.hex())
        self.assertEqual(test_structs[1], obi.decode(encoded, 1))

        # sturct 3
        encoded = obi.encode(test_structs[2], 2)
        self.assertEqual(
            "18a24ec1659380f421195719845bbff4df8568c7bd615812b39436bbcb8ead4c", encoded.hex())
        self.assertEqual(test_structs[2], obi.decode(encoded, 2))

        # sturct 4
        encoded = obi.encode(test_structs[3], 3)
        self.assertEqual(
            "0000000500000003000000010000000200000003000000050000002c00000037000000420000004d00000058000000050000270f000027100000271100002712000027130000000100bc614e00000000", encoded.hex())
        self.assertEqual(test_structs[3], obi.decode(encoded, 3))

        # sturct 5
        encoded = obi.encode(test_structs[4], 4)
        self.assertEqual(
            "0000000300000007fffefd0001020300000003abcdef00000001640000000500000000000000000000000040aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa000000070100010101010000000006000000046d756d75000000046c756c75000000046d6f6d6f00000004746f746f00000004636f636f00000004626f626f", encoded.hex())
        self.assertEqual(test_structs[4], obi.decode(encoded, 4))
