from ..cache_consumer_mock import CACHE_CONSUMER_MOCK
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestCACHE_CONSUMER_MOCK(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(
            CACHE_CONSUMER_MOCK,
            self.test_account1,
            on_install_params={
                # For testing
                "bridge_address": self.test_account1,
                "req_template": bytes.fromhex(
                    "000000047465737400000000000000010000000f00000003425443000000000000006400000000000000040000000000000004"
                ),
            },
        )

    def test_1(self):
        print(">>>> \n ", self.score.get_request_key_template())
