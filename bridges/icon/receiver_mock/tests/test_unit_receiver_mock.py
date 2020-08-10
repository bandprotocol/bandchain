from ..receiver_mock import RECEIVER_MOCK
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestRECEIVER_MOCK(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(RECEIVER_MOCK, self.test_account1)

