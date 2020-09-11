from ..aggregator import Aggregator
from tbears.libs.scoretest.score_test_case import ScoreTestCase
from iconservice.base.exception import IconScoreException


class TestAggregator(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(
            Aggregator, self.test_account1, on_install_params={"bridge_address": self.test_account1}
        )

    def test_get_bridge_address(self):
        self.assertEqual(self.score.get_bridge_address(), self.test_account1)

    def test_symbol_to_indexes(self):
        self.assertEqual(self.score.symbol_to_indexes("YFI"), (1, 0))
        self.assertEqual(self.score.symbol_to_indexes("BCH"), (0, 6))
        self.assertEqual(self.score.symbol_to_indexes("THETA"), (1, 14))
        self.assertEqual(self.score.symbol_to_indexes("SXP"), (2, 24))
        self.assertEqual(self.score.symbol_to_indexes("STORJ"), (5, 24))

        with self.assertRaises(IconScoreException) as e:
            self.score.symbol_to_indexes("AAA")
        self.assertEqual(e.exception.code, 32)
        self.assertEqual(e.exception.message, "UNKNOWN_SYMBOL")
