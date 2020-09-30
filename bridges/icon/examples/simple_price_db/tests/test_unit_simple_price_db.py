from ..simple_price_db import SimplePriceDB
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestSimplePriceDB(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(SimplePriceDB, self.test_account1)

