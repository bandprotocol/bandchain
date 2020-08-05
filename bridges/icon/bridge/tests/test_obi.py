from ..utils import obi
from tbears.libs.scoretest.score_test_case import ScoreTestCase
from iconservice.base.exception import IconScoreException


class TestOBI(ScoreTestCase):

    def setUp(self):
        super().setUp()

    def test_obi(self):
        # BTC, 50, 100
        data = bytes.fromhex("00000003425443000000000000003264")
        symbol, remaining = obi.decode_str(data)
        x, remaining = obi.decode_int(remaining, 64)
        y, remaining = obi.decode_int(remaining, 8)
        self.assertEqual(["BTC", 50, 100], [symbol, x, y])

        # band, 400, 100
        data = bytes.fromhex("0000000462616e64000000000000019064")
        symbol, remaining = obi.decode_str(data)
        x, remaining = obi.decode_int(remaining, 64)
        y, remaining = obi.decode_int(remaining, 8)
        self.assertEqual(["band", 400, 100], [symbol, x, y])

    def test_obi_fail(self):
        def should_fail(data: bytes):
            symbol, remaining = obi.decode_str(data)
            x, remaining = obi.decode_int(remaining, 64)
            y, remaining = obi.decode_int(remaining, 8)

        # fail because exceed the data len
        self.assertRaises(
            IconScoreException,
            should_fail,
            bytes.fromhex("000000034254433200000000000064")
        )

        # fail because decoded int should be at least 1 byte
        self.assertRaises(
            IconScoreException,
            lambda data: obi.decode_int(data, 0),
            bytes.fromhex("00")
        )

        # fail because 9 is not divisible by 8
        self.assertRaises(
            IconScoreException,
            lambda data: obi.decode_int(data, 9),
            bytes.fromhex("0011")
        )
