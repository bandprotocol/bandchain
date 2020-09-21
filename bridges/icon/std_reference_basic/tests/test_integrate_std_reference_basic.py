import os
import json
from iconsdk.builder.call_builder import CallBuilder
from iconsdk.builder.transaction_builder import DeployTransactionBuilder, CallTransactionBuilder
from iconsdk.libs.in_memory_zip import gen_deploy_data_content
from iconsdk.signed_transaction import SignedTransaction
from tbears.libs.icon_integrate_test import IconIntegrateTestBase, SCORE_INSTALL_ADDRESS
from iconservice.base.exception import IconScoreException


DIR_PATH = os.path.abspath(os.path.dirname(__file__))


class TestTest(IconIntegrateTestBase):
    TEST_HTTP_ENDPOINT_URI_V3 = "http://127.0.0.1:9000/api/v3"
    SCORE_PROJECT = os.path.abspath(os.path.join(DIR_PATH, ".."))

    def setUp(self):
        super().setUp()

        self.icon_service = None
        # if you want to send request to network, uncomment next line and set self.TEST_HTTP_ENDPOINT_URI_V3
        # self.icon_service = IconService(HTTPProvider(self.TEST_HTTP_ENDPOINT_URI_V3))

        # install SCORE
        self._std_reference_basic = self._deploy_score()["scoreAddress"]

    def _deploy_score(self, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        # Generates an instance of transaction for deploying SCORE.
        transaction = (
            DeployTransactionBuilder()
            .from_(self._test1.get_address())
            .to(to)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .content_type("application/zip")
            .content(gen_deploy_data_content(self.SCORE_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def set_refs_by_using_relay(self, relay_data):
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_basic)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params({"_json_data_list": json.dumps(relay_data)})
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

    def test_relay_and_get_reference_data(self):
        pairs = ["BTC/ETH", "ETH/BAND", "BAND/ICX", "ICX/SUSD", "SUSD/BTC"]
        relay_data = [
            {"symbol": "BTC", "rate": int(16e9), "resolve_time": 1600419716, "request_id": 1},
            {"symbol": "ETH", "rate": int(8e9), "resolve_time": 1600419717, "request_id": 2},
            {"symbol": "BAND", "rate": int(4e9), "resolve_time": 1600419718, "request_id": 3},
            {"symbol": "ICX", "rate": int(2e9), "resolve_time": 1600419719, "request_id": 4},
            {"symbol": "SUSD", "rate": int(1e9), "resolve_time": 1600419720, "request_id": 5},
        ]
        for pair in pairs:
            with self.assertRaises(IconScoreException) as e:
                call = (
                    CallBuilder()
                    .from_(self._test1.get_address())
                    .to(self._std_reference_basic)
                    .method("get_reference_data")
                    .params({"_pair": pair})
                    .build()
                )
                self.process_call(call, self.icon_service)

            self.assertEqual(e.exception.code, 32)
            self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        self.set_refs_by_using_relay(relay_data)

        expected_rates = [
            (b["rate"] * int(1e18) // q["rate"])
            for (b, q) in list(zip(relay_data, relay_data[1:] + relay_data[:1]))
        ]
        last_update_base_list = [b["resolve_time"] * int(1e6) for b in relay_data]
        last_update_quote_list = [
            b["resolve_time"] * int(1e6) for b in relay_data[1:] + relay_data[:1]
        ]
        for (pair, rate, last_update_base, last_update_quote) in list(
            zip(pairs, expected_rates, last_update_base_list, last_update_quote_list)
        ):
            call = (
                CallBuilder()
                .from_(self._test1.get_address())
                .to(self._std_reference_basic)
                .method("get_reference_data")
                .params({"_pair": pair})
                .build()
            )
            response = self.process_call(call, self.icon_service)
            self.assertEqual(rate, response["rate"])
            self.assertEqual(last_update_base, response["last_update_base"])
            self.assertEqual(last_update_quote, response["last_update_quote"])

    def test_relay_and_get_reference_data_bulk(self):
        pairs = ["BTC/ETH", "ETH/BAND", "BAND/ICX", "ICX/SUSD", "SUSD/BTC"]
        relay_data = [
            {"symbol": "BTC", "rate": int(16e9), "resolve_time": 1600419716, "request_id": 1},
            {"symbol": "ETH", "rate": int(8e9), "resolve_time": 1600419717, "request_id": 2},
            {"symbol": "BAND", "rate": int(4e9), "resolve_time": 1600419718, "request_id": 3},
            {"symbol": "ICX", "rate": int(2e9), "resolve_time": 1600419719, "request_id": 4},
            {"symbol": "SUSD", "rate": int(1e9), "resolve_time": 1600419720, "request_id": 5},
        ]

        with self.assertRaises(IconScoreException) as e:
            call = (
                CallBuilder()
                .from_(self._test1.get_address())
                .to(self._std_reference_basic)
                .method("get_reference_data_bulk")
                .params({"_json_pairs": json.dumps(pairs)})
                .build()
            )
            self.process_call(call, self.icon_service)

        self.assertEqual(e.exception.code, 32)
        self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        self.set_refs_by_using_relay(relay_data)

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_basic)
            .method("get_reference_data_bulk")
            .params({"_json_pairs": json.dumps(pairs)})
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            [
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419716000000,
                    "last_update_quote": 1600419717000000,
                },
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419717000000,
                    "last_update_quote": 1600419718000000,
                },
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419718000000,
                    "last_update_quote": 1600419719000000,
                },
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419719000000,
                    "last_update_quote": 1600419720000000,
                },
                {
                    "rate": 62500000000000000,
                    "last_update_base": 1600419720000000,
                    "last_update_quote": 1600419716000000,
                },
            ],
            response,
        )

