from ..receiver_mock import ReceiverMock
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestReceiverMock(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(ReceiverMock, self.test_account1)

