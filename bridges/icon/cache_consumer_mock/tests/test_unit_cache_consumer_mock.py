from ..cache_consumer_mock import CacheConsunmerMock
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestCacheConsumerMock(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(
            CacheConsunmerMock,
            self.test_account1,
            on_install_params={
                # For testing
                "bridge_address": self.test_account1,
                "req_template": bytes.fromhex(
                    "000000047465737400000000000000010000000f00000003425443000000000000006400000000000000040000000000000004"
                ),
            },
        )

    def test_(self):
        self.assertEqual(
            {
                "client_id": "test",
                "oracle_script_id": 1,
                "calldata": b"\x00\x00\x00\x03BTC\x00\x00\x00\x00\x00\x00\x00d",
                "ask_count": 4,
                "min_count": 4,
            },
            self.score.get_request_key_template(),
        )
