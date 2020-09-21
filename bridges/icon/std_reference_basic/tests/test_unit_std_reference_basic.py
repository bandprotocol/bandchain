from ..std_reference_basic import StdReferenceBasic
from tbears.libs.scoretest.score_test_case import ScoreTestCase
from iconservice.base.exception import IconScoreException


class TestStdReferenceBasic(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(StdReferenceBasic, self.test_account1)

    def _set_test_data(self):
        self.score._set_refs("BTC", 10_953_310_000_000, 1_600_419_436, 1)
        self.score._set_refs("ETH", 386_084_000_000, 1_600_419_436, 2)
        self.score._set_refs("SUSD", 1_003_000_000, 1_600_419_716, 3)

    def test_get_ref_data(self):
        self.assertEqual(
            {"rate": self.score.ONE, "last_update": self.score.block.timestamp, "request_id": 0},
            self.score._get_ref_data("USD"),
        )

        with self.assertRaises(IconScoreException) as e:
            self.score._get_ref_data("BTC")

        self.assertEqual(e.exception.code, 32)
        self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        with self.assertRaises(IconScoreException) as e:
            self.score._get_ref_data("ETH")

        self.assertEqual(e.exception.code, 32)
        self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        with self.assertRaises(IconScoreException) as e:
            self.score._get_ref_data("SUSD")

        self.assertEqual(e.exception.code, 32)
        self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        self._set_test_data()

        self.assertEqual(
            {"rate": 10_953_310_000_000, "last_update": 1_600_419_436_000_000, "request_id": 1},
            self.score._get_ref_data("BTC"),
        )
        self.assertEqual(
            {"rate": 386_084_000_000, "last_update": 1_600_419_436_000_000, "request_id": 2},
            self.score._get_ref_data("ETH"),
        )
        self.assertEqual(
            {"rate": 1_003_000_000, "last_update": 1_600_419_716_000_000, "request_id": 3},
            self.score._get_ref_data("SUSD"),
        )

    def test_get_reference_data(self):
        self._set_test_data()
        self.assertEqual(
            {
                "rate": self.score.ONE * self.score.ONE,
                "last_update_base": self.score.block.timestamp,
                "last_update_quote": self.score.block.timestamp,
            },
            self.score.get_reference_data("USD/USD"),
        )

        self.assertEqual(
            {
                "rate": 10_953_310_000_000 * self.score.ONE,
                "last_update_base": 1_600_419_436_000_000,
                "last_update_quote": self.score.block.timestamp,
            },
            self.score.get_reference_data("BTC/USD"),
        )

        self.assertEqual(
            {
                "rate": (10_953_310_000_000 * self.score.ONE * self.score.ONE) // 386_084_000_000,
                "last_update_base": 1_600_419_436_000_000,
                "last_update_quote": 1_600_419_436_000_000,
            },
            self.score.get_reference_data("BTC/ETH"),
        )

        self.assertEqual(
            {
                "rate": (386_084_000_000 * self.score.ONE * self.score.ONE) // 1_003_000_000,
                "last_update_base": 1_600_419_436_000_000,
                "last_update_quote": 1_600_419_716_000_000,
            },
            self.score.get_reference_data("ETH/SUSD"),
        )

    def test_get_reference_data_bulk(self):
        self._set_test_data()
        self.assertEqual(
            [
                {
                    "rate": self.score.ONE * self.score.ONE,
                    "last_update_base": self.score.block.timestamp,
                    "last_update_quote": self.score.block.timestamp,
                },
                {
                    "rate": (10_953_310_000_000 * self.score.ONE * self.score.ONE)
                    // 386_084_000_000,
                    "last_update_base": 1_600_419_436_000_000,
                    "last_update_quote": 1_600_419_436_000_000,
                },
                {
                    "rate": (386_084_000_000 * self.score.ONE * self.score.ONE) // 1_003_000_000,
                    "last_update_base": 1_600_419_436_000_000,
                    "last_update_quote": 1_600_419_716_000_000,
                },
            ],
            self.score.get_reference_data_bulk('["USD/USD","BTC/ETH","ETH/SUSD"]'),
        )
